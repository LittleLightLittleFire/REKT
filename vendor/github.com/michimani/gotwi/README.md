gotwi
===

[![Go Reference](https://pkg.go.dev/badge/github.com/michimani/gotwi.svg)](https://pkg.go.dev/github.com/michimani/gotwi)
[![Twitter API v2](https://img.shields.io/endpoint?url=https%3A%2F%2Ftwbadges.glitch.me%2Fbadges%2Fv2)](https://developer.twitter.com/en/docs/twitter-api)
[![codecov](https://codecov.io/gh/michimani/gotwi/branch/main/graph/badge.svg?token=NA873TE6RV)](https://codecov.io/gh/michimani/gotwi)

This is a library for using the Twitter API v2 in the Go language. (It is still under development).

# Supported APIs

[Twitter API Documentation | Docs | Twitter Developer Platform  ](https://developer.twitter.com/en/docs/twitter-api)

Progress of supporting APIs:

| Category | Sub Category | Endpoint |
| --- | --- | --- |
| Tweets | Tweet lookup | `GET /2/tweets` |
|  |  | `GET /2/tweets/:id` |
|  | Manage Tweet | `POST /2/tweets` |
|  |  | `DELETE /2/tweets/:id` |
|  | Timelines | `GET /2/users/:id/tweets` |
|  |  | `GET /2/users/:id/mentions` |
|  |  | `GET /2/users/:id/timelines/reverse_chronological` |
|  | Search Tweets | `GET /2/tweets/search/recent` |
|  |  | `GET /2/tweets/search/all` |
|  | Tweet counts | `GET /2/tweets/counts/recent` |
|  |  | `GET /2/tweets/counts/all` |
|  | Filtered stream | `POST /2/tweets/search/stream/rules` |
|  |  | `GET /2/tweets/search/stream/rules` |
|  |  | `GET /2/tweets/search/stream` |
|  | Volume streams | `GET /2/tweets/sample/stream` |
|  | Retweets | `GET /2/users/:id/retweeted_by` |
|  |  | `POST /2/users/:id/retweets` |
|  |  | `DELETE /2/users/:id/retweets/:source_tweet_id` |
|  | Likes | `GET /2/tweets/:id/liking_users` |
|  |  | `GET /2/tweets/:id/liked_tweets` |
|  |  | `POST /2/users/:id/likes` |
|  |  | `DELETE /2/users/:id/likes/:tweet_id` |
|  | Hide replies | `PUT /2/tweets/:id/hidden` |
|  | Quote Tweets | `GET /2/tweets/:id/quote_tweets` |
|  | Bookmarks | `GET /2/users/:id/bookmarks` |
|  |  | `POST /2/users/:id/bookmarks` |
|  |  | `DELETE /2/users/:id/bookmarks/:tweet_id` |
| Users | User lookup | `GET /2/users` |
|  |  | `GET /2/users/:id` |
|  |  | `GET /2/users/by` |
|  |  | `GET /2/users/by/username` |
|  |  | `GET /2/users/by/me` |
|  | Follows | `GET /2/users/:id/following` |
|  |  | `GET /2/users/:id/followers` |
|  |  | `POST /2/users/:id/following` |
|  |  | `DELETE /2/users/:source_user_id/following/:target_user_id` |
|  | Blocks | `GET /2/users/:id/blocking` |
|  |  | `POST /2/users/:id/blocking` |
|  |  | `DELETE /2/users/:source_user_id/blocking/:target_user_id` |
|  | Mutes | `GET /2/users/:id/muting` |
|  |  | `POST /2/users/:id/muting` |
|  |  | `DELETE /2/users/:source_user_id/muting/:target_user_id` |
| Lists | List lookup | `GET /2/lists/:id` |
|  |  | `GET /2/users/:id/owned_lists` |
|  | Manage Lists | `POST /2/lists` |
|  |  | `DELETE /2/lists/:id` |
|  |  | `PUT /2/lists/:id` |
|  | List Tweets lookup | `GET /2/lists/:id/tweets` |
|  | List members | `GET /2/users/:id/list_memberships` |
|  |  | `GET /2/lists/:id/members` |
|  |  | `POST /2/lists/:id/members` |
|  |  | `DELETE /2/lists/:id/members/:user_id` |
|  | List follows | `GET /2/lists/:id/followers` |
|  |  | `GET /2/users/:id/followed_lists` |
|  |  | `POST /2/users/:id/followed_lists` |
|  |  | `DELETE /2/users/:id/followed_lists/:list_id` |
|  | Pinned Lists | `GET /2/users/:id/pinned_lists` |
|  |  | `POST /2/users/:id/pinned_lists` |
|  |  | `DELETE /2/users/:id/pinned_lists/:list_id` |
| Spaces | Spaces Lookup | `GET /2/spaces/:id` |
|  |  | `GET /2/spaces` |
|  |  | `GET /2/spaces/by/creator_ids` |
|  |  | `GET /2/spaces/:id/buyers` |
|  |  | `GET /2/spaces/:id/tweets` |
|  | Search Spaces | `GET /2/spaces/search` |
| Compliance | Batch compliance | `GET /2/compliance/jobs/:id` |
|  |  | `GET /2/compliance/jobs` |
|  |  | `POST /2/compliance/jobs` |


# How to use

## Prepare

Set the API key and API key secret to environment variables.

```bash
export GOTWI_API_KEY='your-api-key'
export GOTWI_API_KEY_SECRET='your-api-key-secret'
```

## Request with OAuth 1.0a User Context

With this authentication method, each operation will be performed as the authenticated Twitter account. For example, you can tweet as that account, or retrieve accounts that are blocked by that account.

### Example: Get your own information.

```go
package main

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"
)

func main() {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           "your-access-token",
		OAuthTokenSecret:     "your-access-token-secret",
	}

	c, err := gotwi.NewClient(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.GetMeInput{
		Expansions: fields.ExpansionList{
			fields.ExpansionPinnedTweetID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
		},
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldCreatedAt,
		},
	}

	u, err := userlookup.GetMe(context.Background(), c, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ID:          ", gotwi.StringValue(u.Data.ID))
	fmt.Println("Name:        ", gotwi.StringValue(u.Data.Name))
	fmt.Println("Username:    ", gotwi.StringValue(u.Data.Username))
	fmt.Println("CreatedAt:   ", u.Data.CreatedAt)
	if u.Includes.Tweets != nil {
		for _, t := range u.Includes.Tweets {
			fmt.Println("PinnedTweet: ", gotwi.StringValue(t.Text))
		}
	}
}
```

```
go run main.go
```

You will get the output like following.

```
ID:           581780917
Name:         michimani Lv.873
Username:     michimani210
CreatedAt:    2012-05-16 12:07:04 +0000 UTC
PinnedTweet:  OpenAI API の Function Calling を使って自然言語で AWS リソースを作成してみる
```

### Example: Tweet with poll.

```go
package main

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

func main() {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           "your-access-token",
		OAuthTokenSecret:     "your-access-token-secret",
	}

	c, err := gotwi.NewClient(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.CreateInput{
		Text: gotwi.String("This is a test tweet with poll."),
		Poll: &types.CreateInputPoll{
			DurationMinutes: gotwi.Int(5),
			Options: []string{
				"Cyan",
				"Magenta",
				"Yellow",
				"Key plate",
			},
		},
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("[%s] %s\n", gotwi.StringValue(res.Data.ID), gotwi.StringValue(res.Data.Text))
}
```

```
go run main.go
```

You will get the output like following.

```
[1462813519607263236] This is a test tweet with poll.
```

## Request with OAuth 2.0 Bearer Token

This authentication method allows only read-only access to public information.

### Example: Get a user by user name.

⚠ This example only works with Twitter API v2 Basic or Pro plan. see details: [Developers Portal](https://developer.twitter.com/en/portal/products)

```go
package main

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"
)

func main() {
	c, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.GetByUsernameInput{
		Username: "michimani210",
		Expansions: fields.ExpansionList{
			fields.ExpansionPinnedTweetID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
		},
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldCreatedAt,
		},
	}

	u, err := userlookup.GetByUsername(context.Background(), c, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ID:          ", gotwi.StringValue(u.Data.ID))
	fmt.Println("Name:        ", gotwi.StringValue(u.Data.Name))
	fmt.Println("Username:    ", gotwi.StringValue(u.Data.Username))
	fmt.Println("CreatedAt:   ", u.Data.CreatedAt)
	if u.Includes.Tweets != nil {
		for _, t := range u.Includes.Tweets {
			fmt.Println("PinnedTweet: ", gotwi.StringValue(t.Text))
		}
	}
}
```

```
go run main.go
```

You will get the output like following.

```
ID:           581780917
Name:         michimani Lv.861
Username:     michimani210
CreatedAt:    2012-05-16 12:07:04 +0000 UTC
PinnedTweet:  真偽をハッキリしたい西城秀樹「ブーリアン、ブーリアン」
```

## Request with OAuth 2.0 Authorization Code with PKCE

If you already have a pre-generated access token (e.g. OAuth 2.0 Authorization Code with PKCE), you can use `NewClientWithAccessToken()` function to generate a Gotwi client.

```go
in := &gotwi.NewClientWithAccessTokenInput{
	AccessToken: "your-access-token",
}

c, err := gotwi.NewClientWithAccessToken(in)
if err != nil {
	// error handling
}
```

See below for information on which authentication methods are available for which endpoints.

[Twitter API v2 authentication mapping | Docs | Twitter Developer Platform  ](https://developer.twitter.com/en/docs/authentication/guides/v2-authentication-mapping)

## Error handling

Each function that calls the Twitter API (e.g. `retweet.ListUsers()`) may return an error for some reason.
If the error is caused by the Twitter API returning a status other than 2XX, you can check the details by doing the following.

```go
res, err := retweet.ListUsers(context.Background(), c, p)
if err != nil {
	fmt.Println(err)

	// more error information
	ge := err.(*gotwi.GotwiError)
	if ge.OnAPI {
		fmt.Println(ge.Title)
		fmt.Println(ge.Detail)
		fmt.Println(ge.Type)
		fmt.Println(ge.Status)
		fmt.Println(ge.StatusCode)

		for _, ae := range ge.APIErrors {
			fmt.Println(ae.Message)
			fmt.Println(ae.Label)
			fmt.Println(ae.Parameters)
			fmt.Println(ae.Code)
			fmt.Println(ae.Code.Detail())
		}

		if ge.RateLimitInfo != nil {
			fmt.Println(ge.RateLimitInfo.Limit)
			fmt.Println(ge.RateLimitInfo.Remaining)
			fmt.Println(ge.RateLimitInfo.ResetAt)
		}
	}
}
```



## More examples

See [_examples](https://github.com/michimani/gotwi/tree/main/_examples) directory.

# Licence

[MIT](https://github.com/michimani/gotwi/blob/main/LICENCE)

# Author

[michimani210](https://twitter.com/michimani210)

