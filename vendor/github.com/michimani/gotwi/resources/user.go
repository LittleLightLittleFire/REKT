package resources

import "time"

type User struct {
	ID                *string            `json:"id"`
	Name              *string            `json:"name"`
	Username          *string            `json:"username"`
	CreatedAt         *time.Time         `json:"created_at,omitempty"`
	Description       *string            `json:"description,omitempty"`
	Entities          *UserEntities      `json:"entities,omitempty"`
	Location          *string            `json:"location,omitempty"`
	PinnedTweetID     *string            `json:"pinned_tweet_id,omitempty"`
	ProfileImageURL   *string            `json:"profile_image_url,omitempty"`
	Protected         *bool              `json:"protected,omitempty"`
	PublicMetrics     *UserPublicMetrics `json:"public_metrics,omitempty"`
	URL               *string            `json:"url,omitempty"`
	Verified          *bool              `json:"verified,omitempty"`
	Withheld          *UserWithheld      `json:"withheld,omitempty"`
	MostRecentTweetID *string            `json:"most_recent_tweet_id,omitempty"`
}

type UserEntities struct {
	URL         *UserURL         `json:"url"`
	Description *UserDescription `json:"description"`
}

type UserURL struct {
	URLs []struct {
		Start       *int    `json:"start"`
		End         *int    `json:"end"`
		URL         *string `json:"url"`
		ExpandedURL *string `json:"expanded_url"`
		DisplayURL  *string `json:"display_url"`
	} `json:"urls"`
}

type UserDescription struct {
	URLs []struct {
		Start       *int    `json:"start"`
		End         *int    `json:"end"`
		URL         *string `json:"url"`
		ExpandedURL *string `json:"expanded_url"`
		DisplayURL  *string `json:"display_url"`
	} `json:"urls"`
	HashTags []UserEntityTag `json:"hashtags"`
	Mentions []UserEntities  `json:"mentions"`
	CashTags []UserEntityTag `json:"cashtags"`
}

type UserEntityTag struct {
	Start *int    `json:"start"`
	End   *int    `json:"end"`
	Tag   *string `json:"tag"`
}

type UserPublicMetrics struct {
	FollowersCount *int `json:"followers_count"`
	FollowingCount *int `json:"following_count"`
	TweetCount     *int `json:"tweet_count"`
	ListedCount    *int `json:"listed_count"`
}

type UserWithheld struct {
	Copyright    *bool     `json:"copyright"`
	CountryCodes []*string `json:"country_codes"`
}
