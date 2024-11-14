package resources

import "time"

type Place struct {
	FullName        *string     `json:"full_name"`
	ID              *string     `json:"id"`
	ContainedWithin *string     `json:"contained_within,omitempty"`
	Country         *string     `json:"country,omitempty"`
	CountryCode     *string     `json:"country_code,omitempty"`
	Geo             *IncludeGeo `json:"geo,omitempty"`
	Name            *string     `json:"name,omitempty"`
	PlaceType       *string     `json:"place_type,omitempty"`
}

type Media struct {
	MediaKey         *string          `json:"media_key"`
	Type             *string          `json:"type"`
	DurationMs       *int             `json:"duration_ms,omitempty"`
	Height           *int             `json:"height,omitempty"`
	NonPublicMetrics map[string]*int  `json:"non_public_metrics,omitempty"`
	OrganicMetrics   map[string]*int  `json:"organic_metrics,omitempty"`
	URL              *string          `json:"url,omitempty"`
	PreviewImageUrl  *string          `json:"preview_image_url,omitempty"`
	PromotedMetrics  map[string]*int  `json:"promoted_metrics,omitempty"`
	PublicMetrics    map[string]*int  `json:"public_metrics,omitempty"`
	Width            *int             `json:"width,omitempty"`
	AltText          *int             `json:"alt_text,omitempty"`
	Variants         []IncludeVariant `json:"variants,omitempty"`
}

type Poll struct {
	ID              *string      `json:"id"`
	Options         []PollOption `json:"options"`
	DurationMinutes *int         `json:"duration_minutes,omitempty"`
	EndDatetime     *time.Time   `json:"end_datetime,omitempty"`
	VotingStatus    *string      `json:"voting_status,omitempty"`
}

type IncludeGeo struct {
	Type       *string     `json:"type"`
	BBox       [4]*float64 `json:"bbox"`
	Properties interface{} `json:"properties"`
}

type PollOption struct {
	Position *int    `json:"position"`
	Label    *string `json:"label"`
	Votes    *int    `json:"votes"`
}

type IncludeVariant struct {
	BitRate     int    `json:"bit_rate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}
