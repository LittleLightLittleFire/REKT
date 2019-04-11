package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "net/http/pprof"

	"golang.org/x/time/rate"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/errwrap"
)

// BotConfig store the bot configuration.
type BotConfig struct {
	BitMexHost            string `json:"bitmex_host"`
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterTokenSecret    string `json:"twitter_token_secret"`
}

func loadConfig() (config BotConfig, err error) {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

// Constants for Websocket
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func runClient(cfg BotConfig, client *twitter.Client, liqChan chan Liquidation) error {
	// Subscribe to the liquidation feed.
	// https://www.bitmex.com/app/wsAPI
	var u url.URL
	u.Scheme = "wss"
	u.Host = cfg.BitMexHost
	u.Path = "realtime"
	u.RawQuery = "subscribe=liquidation"

	// Connect the websocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{})
	if err != nil {
		return errwrap.Wrapf("could not connect to BitMex: {{err}}", err)
	}

	log.Println("Connected to BitMex:", u.String())

	// Handle the pings
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			conn.Close()
		}()

		for range ticker.C {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}()

	// Handle the websocket read
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// The BitMex may "insert" / "delete / "insert" the order when it is able to liquidate at a better price
	// "insert" is sent when the order is submitted
	// "delete" is sent when the order is executed
	// It may also "update" the order when the it is amended or partially filled
	liquidations := make(map[string]Liquidation)

	for {
		var data map[string]interface{}
		if err := conn.ReadJSON(&data); err != nil {
			return err
		}

		if err, ok := data["error"]; ok {
			return fmt.Errorf("error in API response: %v", err)
		}

		// Print JSON so it is parsable later
		rawJSON, _ := json.Marshal(data)
		log.Println(string(rawJSON))

		if table, ok := data["table"]; ok {
			switch table {
			case "liquidation":
				// This will panic if the cast fails, but it is fine, because it meant bitmex sent us bad data
				innerDataList := data["data"].([]interface{})

				switch data["action"] {
				case "partial":
				case "delete":
					for _, innerData := range innerDataList {
						innerData := innerData.(map[string]interface{})
						orderID := innerData["orderID"].(string)

						delete(liquidations, orderID)
					}

				case "update":
					// The liquidation may amended by BitMEX (position may be reduced or price changed)
					for _, innerData := range innerDataList {
						innerData := innerData.(map[string]interface{})

						orderID := innerData["orderID"].(string)

						originalLiq := liquidations[orderID]
						amendedLiq := liquidations[orderID]

						if innerData["price"] != nil {
							amendedLiq.Price = innerData["price"].(float64)
						}

						if innerData["leavesQty"] != nil {
							amendedLiq.Quantity = int64(innerData["leavesQty"].(float64))
						}

						if innerData["symbol"] != nil {
							amendedLiq.Symbol = Symbol(innerData["symbol"].(string))
						}

						if innerData["side"] != nil {
							amendedLiq.Side = innerData["side"].(string)
						}

						difference := amendedLiq.Quantity - originalLiq.Quantity

						// Check if BitMEX is increasing the size if the liquidation order: it means more positions were liquidated
						if difference > 0 {
							// Output a new liquidation based on this difference
							liqChan <- Liquidation{
								PriceQuantity: PriceQuantity{
									Price:    amendedLiq.Price,
									Quantity: difference,
								},
								Symbol:  Symbol(amendedLiq.Symbol),
								Side:    amendedLiq.Side,
								AmendUp: true,
							}
						}

						liquidations[orderID] = amendedLiq
					}

				case "insert":
					for _, innerData := range innerDataList {
						innerData := innerData.(map[string]interface{})

						price := innerData["price"].(float64)
						leavesQty := int64(innerData["leavesQty"].(float64)) // Cast to int64 because this is always int
						symbol := innerData["symbol"].(string)
						side := innerData["side"].(string)
						orderID := innerData["orderID"].(string)

						l := Liquidation{
							PriceQuantity: PriceQuantity{
								Price:    price,
								Quantity: leavesQty,
							},
							Symbol: Symbol(symbol),
							Side:   side,
						}

						liquidations[orderID] = l
						liqChan <- l
					}
				}
			}
		}
	}
}

func symbolLiquidator(state *State, liqChan <-chan Liquidation, tweetChan chan<- string) {
	flusher := time.NewTicker(10 * time.Second)
	defer flusher.Stop()

	var unsentLiquidation *CombinedLiquidation
	var unsentReceivedAt time.Time

	tweet := func(cl CombinedLiquidation) {
		decoration := state.Decorate(cl)
		status := decoration.Apply(cl.String())
		tweetChan <- status
	}

	for {
		select {
		case <-flusher.C:
			if unsentLiquidation == nil {
				continue
			}

			if time.Now().Sub(unsentReceivedAt) < 15*time.Second {
				continue
			}

			// Flush out the current tweet
			tweet(*unsentLiquidation)
			unsentLiquidation = nil

		case l, ok := <-liqChan:
			if !ok {
				return
			}

			log.Println("Got", l)
			if unsentLiquidation == nil {
				combined := l.ToCombined()
				unsentLiquidation = &combined
				unsentReceivedAt = time.Now()
				continue
			}

			// Try and combine
			if unsentLiquidation.CanCombine(l) {
				log.Println("Combining", unsentLiquidation)
				log.Println("With", l)
				unsentLiquidation.Combine(l)
				log.Println("Into", unsentLiquidation)
				continue
			} else {
				log.Println("Can't combine", unsentLiquidation, l)
			}

			// Tweet the existing liquidation if it cannot be combined
			tweet(*unsentLiquidation)
			combined := l.ToCombined()

			unsentLiquidation = &combined
			unsentReceivedAt = time.Now()
		}
	}
}

func liquidator(liqChan <-chan Liquidation, state *State, client *twitter.Client) {
	tweetChan := make(chan string, 10000)
	defer close(tweetChan)
	go func() {
		// https://developer.twitter.com/en/docs/basics/rate-limits
		// 300 tweets in 3 hours -> 100 tweets in 1 hour -> 100 tweets in 3600s
		// -> 36 seconds between every tweet
		limiter := rate.NewLimiter(rate.Every(36*time.Second), 300)
		for status := range tweetChan {
			// Apply the rate limit
			_ = limiter.Wait(context.Background())

			if tweet, _, err := client.Statuses.Update(status, &twitter.StatusUpdateParams{
				TweetMode: "extended",
			}); err != nil {
				log.Println("Failed to tweet:", status, err)
			} else {
				log.Printf("Sent tweet: %v: '%v'\n", tweet.IDStr, status)
			}
		}
	}()

	// Demultiplex this channel by the tickers
	channels := make(map[Symbol]chan Liquidation)
	defer func() {
		for _, c := range channels {
			close(c)
		}
	}()

	for l := range liqChan {
		log.Printf("Detected liqidation: %+v\n", l)
		if channels[l.Symbol] == nil {
			channels[l.Symbol] = make(chan Liquidation, 10000)
			go symbolLiquidator(state, channels[l.Symbol], tweetChan)
		}

		channels[l.Symbol] <- l
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	rand.Seed(time.Now().UnixNano())

	go func() {
		log.Println("Listening on localhost:6060 (pprof)")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Unable to load config:", err)
	}

	state, err := NewState()
	if err != nil {
		log.Fatal("Failed to load state:", err)
	}

	client := twitter.NewClient(oauth1.NewConfig(cfg.TwitterConsumerKey, cfg.TwitterConsumerSecret).Client(oauth1.NoContext, oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterTokenSecret)))
	user, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		log.Fatal("Failed to verify Twitter credentials:", err)
	}

	log.Println("Logged in as:", user.Name)

	// Start the liquidator
	liqChan := make(chan Liquidation, 1024)
	defer close(liqChan)

	go liquidator(liqChan, state, client)

	for {
		if err := runClient(cfg, client, liqChan); err != nil {
			log.Println("Error:", err, "reconnecting in 10 seconds")
			time.Sleep(10 * time.Second)
		}
	}
}
