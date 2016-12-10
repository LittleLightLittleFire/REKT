package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

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

func runClient(cfg BotConfig, twitter *twitter.Client, state *State) error {
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

		for _ = range ticker.C {
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
	// The following sequence is possible: insert ..... update ..... delete/insert ..... update ..... delete/insert ..... delete
	// ..... indicated a posssible time delay

	// Thus we need to keep track of when the order was last deleted and purge it as neccessary
	lastDelete := make(map[string]time.Time)

	for {
		var data map[string]interface{}
		if err := conn.ReadJSON(&data); err != nil {
			return err
		}

		if err, ok := data["error"]; ok {
			return fmt.Errorf("error in API response: %v", err)
		}

		log.Printf("%#v\n", data)

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

						lastDelete[orderID] = time.Now()
					}

				case "update":
					// The liquidation may amended by bitmex (position may be reduced or price changed)

				case "insert":
					for _, innerData := range innerDataList {
						innerData := innerData.(map[string]interface{})

						price := innerData["price"].(float64)
						leavesQty := int64(innerData["leavesQty"].(float64)) // Cast to int64 because this is always int
						symbol := innerData["symbol"].(string)
						side := innerData["side"].(string)
						orderID := innerData["orderID"].(string)

						// Check if this is an insert after a delete
						if _, ok := lastDelete[orderID]; ok {
							continue
						}

						l := Liquidation{
							Price:    price,
							Quantity: leavesQty,
							Symbol:   Symbol(symbol),
							Side:     side,
						}

						dl := state.Decorate(l)
						// TODO: fix this: this does a disk write every time we tweet, which isn't too terrible since we barely do a tweet a second
						if err := state.Save(); err != nil {
							log.Println("Failed to save state:", err)
						}

						status := dl.String()

						if tweet, _, err := twitter.Statuses.Update(status, nil); err != nil {
							log.Println("Failed to tweet:", status)
						} else {
							log.Printf("Sent tweet: %v: '%v'\n", tweet.IDStr, status)
						}
					}
				}
			}
		}

		// Purge expired orders so we don't hemorrhage memory
		now := time.Now()
		for orderID, timestamp := range lastDelete {
			if now.Sub(timestamp) > 10*time.Second {
				delete(lastDelete, orderID)
			}
		}
	}

}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	rand.Seed(time.Now().UnixNano())

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

	if err := runClient(cfg, client, state); err != nil {
		log.Fatal("Error:", err)
	}
}
