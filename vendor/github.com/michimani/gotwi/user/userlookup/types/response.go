package types

import "github.com/michimani/gotwi/resources"

// ListOutput is struct for response of `GET /2/users`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users
type ListOutput struct {
	Data     []resources.User `json:"data"`
	Includes struct {
		Tweets []resources.Tweet `json:"tweets"`
	} `json:"includes"`
	Errors []resources.PartialError `json:"errors"`
}

func (r *ListOutput) HasPartialError() bool {
	return !(r.Errors == nil || len(r.Errors) == 0)
}

// GetOutput is struct for response of `GET /2/users/:id`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-id
type GetOutput struct {
	Data     resources.User `json:"data"`
	Includes struct {
		Tweets []resources.Tweet `json:"tweets"`
	} `json:"includes"`
	Errors []resources.PartialError `json:"errors"`
}

func (r *GetOutput) HasPartialError() bool {
	return !(r.Errors == nil || len(r.Errors) == 0)
}

// ListByUsernamesOutput is struct for response of `GET /2/users/by`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by
type ListByUsernamesOutput struct {
	Data     []resources.User `json:"data"`
	Includes struct {
		Tweets []resources.Tweet `json:"tweets"`
	} `json:"includes"`
	Errors []resources.PartialError `json:"errors"`
}

func (r *ListByUsernamesOutput) HasPartialError() bool {
	return !(r.Errors == nil || len(r.Errors) == 0)
}

// GetByUsernameOutput is struct for response of `GET /2/users/by/username/:username`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by-username-username
type GetByUsernameOutput struct {
	Data     resources.User `json:"data"`
	Includes struct {
		Tweets []resources.Tweet `json:"tweets"`
	} `json:"includes"`
	Errors []resources.PartialError `json:"errors"`
}

func (r *GetByUsernameOutput) HasPartialError() bool {
	return !(r.Errors == nil || len(r.Errors) == 0)
}

// GetMeOutput is struct for response of `GET /2/users/me`.
// more information: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me
type GetMeOutput struct {
	Data     resources.User `json:"data"`
	Includes struct {
		Tweets []resources.Tweet `json:"tweets"`
	} `json:"includes"`
	Errors []resources.PartialError `json:"errors"`
}

func (r *GetMeOutput) HasPartialError() bool {
	return !(r.Errors == nil || len(r.Errors) == 0)
}
