package types

import (
	"io"
	"net/url"
	"strings"

	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/internal/util"
)

// ListInput is struct for requesting `GET /2/users`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users
type ListInput struct {
	accessToken string

	// Query parameters
	IDs         []string // required
	Expansions  fields.ExpansionList
	TweetFields fields.TweetFieldList
	UserFields  fields.UserFieldList
}

var listQueryParameters = map[string]struct{}{
	"ids":          {},
	"expansions":   {},
	"tweet.fields": {},
	"user.fields":  {},
}

func (p *ListInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *ListInput) AccessToken() string {
	return p.accessToken
}

func (p *ListInput) ResolveEndpoint(endpointBase string) string {
	endpoint := endpointBase

	if p.IDs == nil || len(p.IDs) == 0 {
		return ""
	}

	pm := p.ParameterMap()
	if len(pm) > 0 {
		qs := util.QueryString(pm, listQueryParameters)
		endpoint += "?" + qs
	}

	return endpoint
}

func (p *ListInput) Body() (io.Reader, error) {
	return nil, nil
}

func (p *ListInput) ParameterMap() map[string]string {
	m := map[string]string{}

	m["ids"] = util.QueryValue(p.IDs)

	m = fields.SetFieldsParams(m, p.Expansions, p.TweetFields, p.UserFields)

	return m
}

// GetInput is struct for requesting `GET /2/users/:id`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-id
type GetInput struct {
	accessToken string

	// Path parameters
	ID string // required

	// Query parameters
	Expansions  fields.ExpansionList
	TweetFields fields.TweetFieldList
	UserFields  fields.UserFieldList
}

var getQueryParameters = map[string]struct{}{
	"expansions":   {},
	"tweet.fields": {},
	"user.fields":  {},
}

func (p *GetInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *GetInput) AccessToken() string {
	return p.accessToken
}

func (p *GetInput) ResolveEndpoint(endpointBase string) string {
	if p.ID == "" {
		return ""
	}

	encoded := url.QueryEscape(p.ID)
	endpoint := strings.Replace(endpointBase, ":id", encoded, 1)

	pm := p.ParameterMap()
	if len(pm) > 0 {
		qs := util.QueryString(pm, getQueryParameters)
		endpoint += "?" + qs
	}

	return endpoint
}

func (p *GetInput) Body() (io.Reader, error) {
	return nil, nil
}

func (p *GetInput) ParameterMap() map[string]string {
	m := map[string]string{}

	m = fields.SetFieldsParams(m, p.Expansions, p.TweetFields, p.UserFields)

	return m
}

// ListByUsernamesInput is struct for requesting `GET /2/users/by`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by
type ListByUsernamesInput struct {
	accessToken string

	// Query parameters
	Usernames   []string // required
	Expansions  fields.ExpansionList
	TweetFields fields.TweetFieldList
	UserFields  fields.UserFieldList
}

var listByUsernamesQueryParameters = map[string]struct{}{
	"usernames":    {},
	"expansions":   {},
	"tweet.fields": {},
	"user.fields":  {},
}

func (p *ListByUsernamesInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *ListByUsernamesInput) AccessToken() string {
	return p.accessToken
}

func (p *ListByUsernamesInput) ResolveEndpoint(endpointBase string) string {
	endpoint := endpointBase

	if p.Usernames == nil || len(p.Usernames) == 0 {
		return ""
	}

	pm := p.ParameterMap()
	if len(pm) > 0 {
		qs := util.QueryString(pm, listByUsernamesQueryParameters)
		endpoint += "?" + qs
	}

	return endpoint
}

func (p *ListByUsernamesInput) Body() (io.Reader, error) {
	return nil, nil
}

func (p *ListByUsernamesInput) ParameterMap() map[string]string {
	m := map[string]string{}

	m["usernames"] = util.QueryValue(p.Usernames)

	m = fields.SetFieldsParams(m, p.Expansions, p.TweetFields, p.UserFields)

	return m
}

// GetByUsernameInput is struct for requesting `GET /2/users/by/username/:username`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by-username-username
type GetByUsernameInput struct {
	accessToken string

	// Path parameters
	Username string // required

	// Query parameters
	Expansions  fields.ExpansionList
	TweetFields fields.TweetFieldList
	UserFields  fields.UserFieldList
}

var getByUsernameQueryParameters = map[string]struct{}{
	"expansions":   {},
	"tweet.fields": {},
	"user.fields":  {},
}

func (p *GetByUsernameInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *GetByUsernameInput) AccessToken() string {
	return p.accessToken
}

func (p *GetByUsernameInput) ResolveEndpoint(endpointBase string) string {
	if p.Username == "" {
		return ""
	}

	encoded := url.QueryEscape(p.Username)
	endpoint := strings.Replace(endpointBase, ":username", encoded, 1)

	pm := p.ParameterMap()
	if len(pm) > 0 {
		qs := util.QueryString(pm, getByUsernameQueryParameters)
		endpoint += "?" + qs
	}

	return endpoint
}

func (p *GetByUsernameInput) Body() (io.Reader, error) {
	return nil, nil
}

func (p *GetByUsernameInput) ParameterMap() map[string]string {
	m := map[string]string{}
	m = fields.SetFieldsParams(m, p.Expansions, p.TweetFields, p.UserFields)

	return m
}

// GetMeInput is struct for requesting `GET /2/users/me`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me
type GetMeInput struct {
	accessToken string

	// Query parameters
	Expansions  fields.ExpansionList
	TweetFields fields.TweetFieldList
	UserFields  fields.UserFieldList
}

var getMeQueryParameters = map[string]struct{}{
	"expansions":   {},
	"tweet.fields": {},
	"user.fields":  {},
}

func (p *GetMeInput) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *GetMeInput) AccessToken() string {
	return p.accessToken
}

func (p *GetMeInput) ResolveEndpoint(endpointBase string) string {
	endpoint := endpointBase

	pm := p.ParameterMap()
	if len(pm) > 0 {
		qs := util.QueryString(pm, getMeQueryParameters)
		endpoint += "?" + qs
	}

	return endpoint
}

func (p *GetMeInput) Body() (io.Reader, error) {
	return nil, nil
}

func (p *GetMeInput) ParameterMap() map[string]string {
	m := map[string]string{}
	m = fields.SetFieldsParams(m, p.Expansions, p.TweetFields, p.UserFields)
	return m
}
