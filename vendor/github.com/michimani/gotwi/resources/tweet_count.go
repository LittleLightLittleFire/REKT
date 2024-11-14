package resources

import "time"

type TweetCount struct {
	End        *time.Time `json:"end"`
	Start      *time.Time `json:"start"`
	TweetCount *int       `json:"tweet_count"`
}
