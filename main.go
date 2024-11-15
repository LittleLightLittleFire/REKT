package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "net/http/pprof"

	"golang.org/x/time/rate"

	"github.com/gorilla/websocket"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	ctypes "github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/michimani/gotwi/user/userlookup"
	utypes "github.com/michimani/gotwi/user/userlookup/types"
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
		return fmt.Errorf("could not connect to BitMex: %w", err)
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

			case "update":
				var update []Instrument
				if err := json.Unmarshal(data.Data, &update); err != nil {
					return err
				}

				for _, v := range update {
					it.Update(v)
				}
			}

		case "liquidation":
			log.Printf("Received: %v %v %v\n", data.Table, data.Action, string(data.Data))

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
				var update []RawLiquidation
				if err := json.Unmarshal(data.Data, &update); err != nil {
					return err
				}

				// Update last seen to keep it alive
				for _, v := range update {
					lastSeen[v.OrderID] = time.Now()
				}

			case "insert":
				// Wait for instruments table to be loaded
				if it == nil {
					continue
				}

				// Prune last seen
				for k, v := range lastSeen {
					if time.Since(v) > 24*time.Hour {
						delete(lastSeen, k)
					}
				}

				var inserts []RawLiquidation
				if err := json.Unmarshal(data.Data, &inserts); err != nil {
					return err
				}

				for _, v := range inserts {
					if _, ok := lastSeen[v.OrderID]; ok {
						lastSeen[v.OrderID] = time.Now()
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

			if time.Since(unsentReceivedAt) < unsentCombiningDelay {
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

func liquidator(liqChan <-chan Liquidation, state *State, client *gotwi.Client) {
	tweetChan := make(chan preparedTweet, 10000)
	defer close(tweetChan)
	go func() {
		// https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/api-reference/post-tweets
		// 200 requests in 15 min
		// 1500 tweets per 30 days on free plan (50 daily)
		//
		// We don't want to blow all of it in a single day, cap at max 50 tweets a day and refil at 1 every 1728s (1500 / 30 days)
		limiter := rate.NewLimiter(rate.Every(1728*time.Second), 50)

		for status := range tweetChan {
			var minValue float64
			switch {
			case limiter.Burst() < 5:
				minValue = 5000000
			case limiter.Burst() < 10:
				minValue = 1000000
			case limiter.Burst() < 25:
				minValue = 100000
			}

			if status.usdValue < minValue {
				log.Printf("Tweet dropped because of value cap: %v < %v\n", status.usdValue, minValue)
				continue
			}

			// Apply the rate limit
			_ = limiter.Wait(context.Background())

			lag := time.Since(status.timestamp)
			if client != nil {
				res, err := managetweet.Create(context.Background(), client, &ctypes.CreateInput{
					Text: gotwi.String(status.status),
				})
				if err != nil {
					log.Println("Failed to tweet:", status.status, err)
				} else {
					if res.Data.ID != nil {
						log.Printf("Sent tweet: %v: bursts %v: lag %v: '%v'\n", *res.Data.ID, limiter.Burst(), lag, status.status)
					} else {
						log.Printf("Sent tweet: (???): bursts %v: lag %v: '%v'\n", limiter.Burst(), lag, status.status)
					}
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

	go func() {
		log.Println("Listening on localhost:6060 (pprof)")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Unable to load config:", err)
	}

	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               cfg.TwitterConsumerKey,
		APIKeySecret:         cfg.TwitterConsumerSecret,
		OAuthToken:           cfg.TwitterAccessToken,
		OAuthTokenSecret:     cfg.TwitterTokenSecret,
	}

	var client *gotwi.Client
	if cfg.TwitterConsumerKey != "" {
		client, err = gotwi.NewClient(in)
		if err != nil {
			log.Fatalln(err)
		}

		u, err := userlookup.GetMe(context.Background(), client, &utypes.GetMeInput{})
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Logged in as:", *u.Data.Username)
	}

	state, err := NewState()
	if err != nil {
		log.Fatalln("Failed to load state:", err)
	}

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
