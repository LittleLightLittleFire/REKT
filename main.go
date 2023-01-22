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
	"strings"
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

func runClient(cfg BotConfig, liqChan chan Liquidation) error {
	// Subscribe to the liquidation feed.
	// https://www.bitmex.com/app/wsAPI
	var u url.URL
	u.Scheme = "wss"
	u.Host = cfg.BitMexHost
	u.Path = "realtime"
	u.RawQuery = "subscribe=instrument,liquidation"

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

	// Prevent orderIDs from appearing twice.
	lastSeen := make(map[string]time.Time)

	var it *InstrumentTable

	for {
		var data struct {
			Table  string          `json:"table"`
			Action string          `json:"action"`
			Error  string          `json:"error"`
			Data   json.RawMessage `json:"data"`
		}
		if err := conn.ReadJSON(&data); err != nil {
			return err
		}

		log.Printf("Received: %v %v %v\n", data.Table, data.Action, string(data.Data))

		if data.Error != "" {
			return fmt.Errorf("error in API response: %v", err)
		}

		switch data.Table {
		case "instrument":
			switch data.Action {
			case "partial":
				var curr []Instrument
				if err := json.Unmarshal(data.Data, &curr); err != nil {
					return err
				}

				it = NewInstrumentTable(curr)
			}

		case "liquidation":

			// BitMex may "insert" / "delete / "insert" the order when it is able to liquidate at a better price
			// "insert" is sent when the order is submitted
			// "delete" is sent when the order is executed
			// It may also "update" the order when the it is amended or partially filled

			switch data.Action {
			case "partial":
				var curr []RawLiquidation
				if err := json.Unmarshal(data.Data, &curr); err != nil {
					return err
				}

				// Load the current liquidations as last seen
				for _, v := range curr {
					lastSeen[v.OrderID] = time.Now()
				}

			case "update":
				// Ignored, since once the tweet goes out there is no recovering it

			case "delete":
				// Ignored

			case "insert":
				// Wait for instruments table to be loaded
				if it == nil {
					continue
				}

				// Prune last seen
				for k, v := range lastSeen {
					if time.Now().Sub(v) > 4*time.Hour {
						delete(lastSeen, k)
					}
				}

				var inserts []RawLiquidation
				if err := json.Unmarshal(data.Data, &inserts); err != nil {
					return err
				}

				for _, v := range inserts {
					if _, ok := lastSeen[v.OrderID]; ok {
						continue
					}

					lastSeen[v.OrderID] = time.Now()

					l, err := it.Process(v)
					if err != nil {
						log.Printf("failed to process: %+v %v\n", v, err)
						continue
					}

					liqChan <- l
				}
			}
		}
	}
}

func symbolLiquidator(state *State, liqChan <-chan Liquidation, tweetChan chan<- preparedTweet) {
	flusher := time.NewTicker(10 * time.Second)
	defer flusher.Stop()

	var unsentLiquidation *CombinedLiquidation
	var unsentReceivedAt time.Time
	var unsentCombiningDelay time.Duration

	tweet := func(cl CombinedLiquidation) {
		decoration := state.Decorate(cl)
		status := decoration.Apply(cl.String())
		tweetChan <- preparedTweet{
			timestamp: time.Now(),
			usdValue:  cl.USDValue(),
			status:    status,
		}
	}

	newUnsent := func(l Liquidation) {
		combined := l.ToCombined()

		unsentLiquidation = &combined
		unsentReceivedAt = time.Now()
		unsentCombiningDelay = l.CombiningDelay()
	}

	for {
		select {
		case <-flusher.C:
			if unsentLiquidation == nil {
				continue
			}

			if time.Now().Sub(unsentReceivedAt) < unsentCombiningDelay {
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
				newUnsent(l)
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
			newUnsent(l)
		}
	}
}

type preparedTweet struct {
	timestamp time.Time
	usdValue  float64
	status    string
}

func liquidator(liqChan <-chan Liquidation, state *State, client *twitter.Client) {
	tweetChan := make(chan preparedTweet, 10000)
	defer close(tweetChan)
	go func() {
		// https://developer.twitter.com/en/docs/basics/rate-limits
		// 300 tweets in 3 hours -> 100 tweets in 1 hour -> 100 tweets in 3600s
		// -> 36 seconds between every tweet
		limiter := rate.NewLimiter(rate.Every(36*time.Second), 300)

		var lagMode bool

		for status := range tweetChan {
			lag := time.Now().Sub(status.timestamp)
			if lag > 3*time.Minute {
				if !lagMode {
					log.Println("Lag mode enabled", status.timestamp, "more than 3 minutes behind")
				}
				lagMode = true
			} else if lag < 36*time.Second {
				if lagMode {
					log.Println("Lag mode disabled, tweet channel cleared")
				}
				lagMode = false
			}

			if lagMode {
				if status.usdValue < 1000000 {
					log.Printf("Tweet dropped because of lag mode: %+v\n", status)
					continue
				}
			}

			// Apply the rate limit
			_ = limiter.Wait(context.Background())

			if client != nil {
				if tweet, _, err := client.Statuses.Update(status.status, &twitter.StatusUpdateParams{
					TweetMode: "extended",
				}); err != nil {
					log.Println("Failed to tweet:", status.status, err)
					if strings.Contains(err.Error(), "User is over daily status update limit") {
						log.Println("Daily status update limit exceeded, forcing 3m sleep")
						// Force lag mode to activate.
						time.Sleep(3 * time.Minute)
					}
				} else {
					log.Printf("Sent tweet: %v: lag %v: '%v'\n", tweet.IDStr, lag, status.status)
				}
			} else {
				log.Printf("Would have tweeted: lag %v: '%v'\n", lag, status.status)
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
		log.Printf("Detected liquidation: %+v\n", l)

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
		if err := runClient(cfg, liqChan); err != nil {
			log.Println("Error:", err, "reconnecting in 10 seconds")
			time.Sleep(10 * time.Second)
		}
	}
}
