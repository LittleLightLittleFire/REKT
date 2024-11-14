package types

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

// CreateInput is struct for the parameters
// that used for calling POST /2/tweets API.
type CreateInput struct {
	accessToken string

	// JSON body parameter
	DirectMessageDeepLink *string           `json:"direct_message_deep_link,omitempty"`
	ForSuperFollowersOnly *bool             `json:"for_super_followers_only,omitempty"`
	Geo                   *CreateInputGeo   `json:"geo,omitempty"`
	Media                 *CreateInputMedia `json:"media,omitempty"`
	Poll                  *CreateInputPoll  `json:"poll,omitempty"`
	QuoteTweetID          *string           `json:"quote_tweet_id,omitempty"`
	Reply                 *CreateInputReply `json:"reply,omitempty"`
	ReplySettings         *string           `json:"reply_settings,omitempty"`
	Text                  *string           `json:"text,omitempty"`
}

type CreateInputGeo struct {
	PlaceID *string `json:"place_id,omitempty"`
}

type CreateInputMedia struct {
	MediaIDs     []string `json:"media_ids,omitempty"`
	TaggedUserID *string  `json:"tagged_user_ids,omitempty"`
}

type CreateInputPoll struct {
	DurationMinutes *int     `json:"duration_minutes,omitempty"`
	Options         []string `json:"options,omitempty"`
}

type CreateInputReply struct {
	ExcludeReplyUserIDs []string `json:"exclude_reply_user_ids,omitempty"`
	InReplyToTweetID    string   `json:"in_reply_to_tweet_id,omitempty"`
}

func (p *CreateInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *CreateInput) AccessToken() string {
	return p.accessToken
}

func (p *CreateInput) ResolveEndpoint(endpointBase string) string {
	return endpointBase
}

func (p *CreateInput) Body() (io.Reader, error) {
	json, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(json)), nil
}

func (p *CreateInput) ParameterMap() map[string]string {
	return map[string]string{}
}

type DeleteInput struct {
	accessToken string

	// Path parameter
	ID string // required: The tweet ID to delete
}

func (p *DeleteInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *DeleteInput) AccessToken() string {
	return p.accessToken
}

func (p *DeleteInput) ResolveEndpoint(endpointBase string) string {
	if p.ID == "" {
		return ""
	}

	escaped := url.QueryEscape(p.ID)
	endpoint := strings.Replace(endpointBase, ":id", escaped, 1)

	return endpoint
}

func (p *DeleteInput) Body() (io.Reader, error) {
	return nil, nil
}

func (p *DeleteInput) ParameterMap() map[string]string {
	return map[string]string{}
}
