CHANGELOG 
===
This is a version of CHANGELOG less than v1.0.0

## [Unreleased]

* TBD

v0.16.1 (2024-09-24)
===

### Features

* Support `MostRecentTweetID` fields in `User` struct. ([cf40a01](https://github.com/michimani/gotwi/commit/cf40a01b7936d472565382de183d6b167fd29ce2) by [@AtakanPehlivanoglu](https://github.com/AtakanPehlivanoglu))

v0.16.0 (2024-09-19)
===

### Features

* Ability to set API Key & API Key Secret manually. ([8abb233](https://github.com/michimani/gotwi/commit/8abb23320f5b912392e397d89689809252fd6d36) by [@EricFrancis12](https://github.com/EricFrancis12))

v0.15.0 (2024-06-25)
===

### Features

* Support `NoteTweet` fields in `Tweet` struct. ([370f3eb](https://github.com/michimani/gotwi/commit/370f3ebda03bac61d4b3fbfdb3bb8dc4d53d2bc8) by [@NHypocrite](https://github.com/NHypocrite))

### Fixes
* bump Go version 1.22

v0.14.0 (2023-06-14)
===

### Fixes
* bump Go version 1.20

v0.13.0 (2022-11-12)
===

### ⚠ BREAKING CHANGES

**This version is not compatible with v0.12.5 or earlier.**

**If your application uses the `PartialError.ResourceId` and `Tweet.ConversationId` fields, you will need to modify them when updating to this version. `PartialError.ResourceID` and `Tweet.ConversationID` instead of them, respectively.**

### Fixes

* Fix invalid field names. ([635bd32](https://github.com/michimani/gotwi/commit/635bd326a8162aea88a39324283630a416ecc0cc))

v0.12.5 (2022-11-12)
===

### Features

* Support `EditHistoryTweetIDs` fields in `Tweet` struct and support `MatchingRules` for response of `GET /2/tweets/search/stream`. ([bd569c4](https://github.com/michimani/gotwi/commit/bd569c412e9d3e7e6822e798664344e6fd93e917))

v0.12.4 (2022-10-30)
===

### Features

* Support all fields in `includes` field for response of `GET /2/tweets/search/recent` and `GET /2/tweets/search/all`. ([bc396a3](https://github.com/michimani/gotwi/commit/bc396a33264407399b4fd56a21f1e1c87d310b74))

v0.12.3 (2022-10-20)
===

### Features

* Support `media.variants` field for response of `GET /2/tweets/search/stream` and so on. ([b453dc2](https://github.com/michimani/gotwi/commit/b453dc2449c3d3f56395831df0356bb9666246b9))

v0.12.2 (2022-08-29)
===

### Fixes

* Fix struct of response for `GET /2/users/:id/tweets`, `GET /2/users/:id/mentions` and `GET /2/users/:id/timelines/reverse_chronological`. ([25c8bf0](https://github.com/michimani/gotwi/commit/25c8bf055dc05c94c499b3d2856109c1567c370b)), ([#288](https://github.com/michimani/gotwi/issues/288)) ([56f1c42](https://github.com/michimani/gotwi/commit/56f1c420a7a2ec1333662941dbf0a5fa0fd12669))

v0.12.1 (2022-08-17)
===

### Fixes

* Fix struct of response for `PUT /2/lists/:id`. ([#284](https://github.com/michimani/gotwi/issues/284)) ([f131485](https://github.com/michimani/gotwi/commit/f131485efb2d671b12d6c424d1611042630aa5ca))

v0.12.0 (2022-08-04)
===

### New Supported APIs

* Support New API: `GET /2/users/:id/timelines/reverse_chronological` ([#268](https://github.com/michimani/gotwi/issues/268)) ([f8153fa](https://github.com/michimani/gotwi/commit/f8153fa9907edaab0ad54a6f0bd4e0e3d8036cac))

### Features

* Support `max_results` and `pagination_token` parameter at `GET /2/tweets/:id/liking_users`. ([#272](https://github.com/michimani/gotwi/issues/272)) ([308e2b8](https://github.com/michimani/gotwi/commit/308e2b8ba287d98a934b33f56fdaee90fa741b29))
* Support `max_results` and `pagination_token` parameter at `GET /2/tweets/:id/retweeted_by`. ([#270](https://github.com/michimani/gotwi/issues/270)) ([279e364](https://github.com/michimani/gotwi/commit/279e364b64362d72192dd27385666793b5488e4c))


### Fixes

* Fix invalid field type. ([#279](https://github.com/michimani/gotwi/issues/279)) ([f8153fa](https://github.com/michimani/gotwi/commit/f8153fa9907edaab0ad54a6f0bd4e0e3d8036cac))
* Remove invalid input fields. ([#272](https://github.com/michimani/gotwi/issues/272)) ([308e2b8](https://github.com/michimani/gotwi/commit/308e2b8ba287d98a934b33f56fdaee90fa741b29)) ([#270](https://github.com/michimani/gotwi/issues/270)) ([279e364](https://github.com/michimani/gotwi/commit/279e364b64362d72192dd27385666793b5488e4c))

v0.11.8 (2022-08-04)
===

### Features

* Support `quote_tweet_id` parameter at `POST /2/tweets`. ([#267](https://github.com/michimani/gotwi/issues/267)) ([b6de378](https://github.com/michimani/gotwi/commit/b6de37861aaa383f08697593d0c1c07a8cf8954f))
* Support `backfill_minutes` parameter at `GET /2/tweets/sample/stream`. ([#269](https://github.com/michimani/gotwi/issues/269)) ([b77035d](https://github.com/michimani/gotwi/commit/b77035d44d93b63c216e43dc6e7671fe06a140c0))
* Support `exclude` parameter at `GET /2/tweets/:id/quote_tweets`. ([#271](https://github.com/michimani/gotwi/issues/271)) [6b319d2](https://github.com/michimani/gotwi/commit/6b319d2908d676548eb9e7238a9f485edf365690)

v0.11.7 (2022-08-03)
===

### Features

* Support `sort_order` parameter for Search Tweet API. ([#263](https://github.com/michimani/gotwi/issues/263)) ([193f336](https://github.com/michimani/gotwi/commit/193f3362c664844a3d11cbc53fb56728ec51e96f))

v0.11.6 (2022-07-06)
===

### Fixes

* `filteredstream.SearchStream` sometimes fails with "unexpected end of JSON input" error. ([#259](https://github.com/michimani/gotwi/issues/259))

v0.11.5 (2022-06-12)
===

### Fixes

* TweetEntities struct has some invalid and not enough fields. ([#256](https://github.com/michimani/gotwi/pull/256/files))

v0.11.4 (2022-05-27)
===

### Fixes

* Media struct does not have URL field. ([#252](https://github.com/michimani/gotwi/pull/252/files))

v0.11.3 (2022-05-27)
===

### Features

* Add debug mode.

### Fixes

* Tweet lookup API response does not include Tweets, Places, Media, and Polls Fields. ([#248](https://github.com/michimani/gotwi/pull/248/files))

v0.11.2 (2022-04-21)
===

### New Supported APIs

* `POST /2/tweets/search/stream/rules`
* `GET /2/tweets/search/stream`

v0.11.1 (2022-04-19)
===

### New Supported APIs

* `GET /2/tweets/sample/stream`

v0.11.0 (2022-04-13)
===

### ⚠ BREAKING CHANGES

**This version is not compatible with v0.10.4 or earlier. The structure of the library has changed significantly in `v0.11.0` and later compared to `v0.10.4` and earlier. If you are using `v0.10.4` or earlier, updating to `v0.11.0` or later may result in your application not working properly.**

* Change package names, function names, and method names to more descriptive names. 

v0.10.4 (2022-04-12)
===

### Documentation

* Update unreleased comment.

v0.10.3 (2022-04-10)
===

### New Supported APIs
* `GET /2/compliance/jobs/:id`
* `GET /2/compliance/jobs`
* `POST /2/compliance/jobs`

v0.10.2 (2022-04-06)
===

### New Supported APIs
* `GET /2/users/:id/bookmarks`
* `POST /2/users/:id/bookmarks`
* `DELETE /2/users/:id/bookmarks/:tweet_id`

### Features
* Generate GotwiClient with access token

v0.10.1 (2022-03-31)
===

### New Supported APIs
* `GET /2/tweets/:id/quote_tweets`

### Fixes
* bump Go version 1.18

v0.10.0 (2022-02-01)
===

### Features
* Handling API errors


v0.9.10 (2022-01-14)
===

### New Supported APIs
* `GET /2/spaces/:id/buyers`
* `GET /2/spaces/:id/tweets`

v0.9.9 (2022-01-14)
===

### Documentation
* add some tests

### Fixes
* use `io.ReadAll` instead of `ioutil.ReadAll` 
* comment for API (`users.UserLookupMe`)

v0.9.8 (2021-12-28)
===

### New Supported APIs
* `GET /2/users/me`

### Documentation
* add some tests

v0.9.7 (2021-12-24)
===

### Fixes
* remove some unnecessary processing

### Documentation
* add some tests

v0.9.6 (2021-12-24)
===

### Documentation
* add code coverage tool

v0.9.5 (2021-12-01)
===

### Fixes
* Creating OAuth 1.0 signature
* some tests

v0.9.4 (2021-11-22)
===

### Documentation
* add examples
* add pkg.go badge

v0.9.3 (2021-11-22)
===

### New Supported APIs
* `DELETE /2/tweets`
* `POST /2/tweets`

v0.9.2 (2021-11-22)
===

### New Supported APIs
* `GET /2/users/:id/list_memberships`
* `GET /2/users/:id/followed_lists`
* `GET /2/lists/:id/followers`
* `GET /2/users/:id/pinned_lists`
* `GET /2/lists/:id/tweets`

v0.9.1 (2021-11-18)
===

### New Supported APIs
* `GET /2/lists/:id/members`

v0.9.0 (2021-11-18)
===

### New Supported APIs
* `GET /2/users/:id/owned_lists`
* `GET /2/lists/:id`

### Fixes
* type of resource fields
* name of Lists resources


v0.8.2 (2021-11-16)
===

### New Supported APIs
* `GET /2/tweets/search/stream/rules`

v0.8.1 (2021-11-12)
===

### New Supported APIs
* `GET /2/spaces/by/creator_ids`
* `GET /2/spaces`
* `GET /2/spaces/:id`
* `GET /2/spaces/search`

### Fixes
* some fields

v0.8.0 (2021-10-29)
===

### Features
* Call API with context

### Fixes
* Use json decoder

v0.7.0 (2021-10-27)
===

### Fixes
* type of some fields

v0.6.0 (2021-10-27)
===

### New Supported APIs
* `POST DELETE /2/users/:id/pinned_lists`
* `POST DELETE /2/lists/:id/follows`
* `POST DELETE /2/lists/:id/members`
* `DELETE /2/lists/:id`
* `PUT /2/lists/:id`
* `POST /2/lists`
* `PUT /2/tweets/:id/hidden`

### Fixes
* type of JSON parameters
* not ok error struct

v0.5.2 (2021-10-21)
===

### New Supported APIs
* `DELETE /2/users/:id/retweets/:source_tweet_id`
* `POST /2/users/:id/retweets`
* `DELETE /2/users/:id/likes`
* `POST /2/users/:id/likes`
* `GET /2/tweets/:id/liked_tweets`
* `DELETE /2/users/:source_user_id/muting/:target_user_id`
* `POST /2/users/:id/muting`
* `DELETE /2/users/:source_user_id/blocking/:target_user_id`
* `POST /2/users/:id/blocking`
* `DELETE /2/users/:source_user_id/following/:target_user_id`
* `POST /2/users/:id/following`
* `GET /2/tweets/:id/liking_users`
* `GET /2/users/:id/retweeted_by`
* `GET /2/users/:id/muting`
* `GET /2/users/:id/mentions`
* `GET /2/users/:id/tweets`

v0.5.1 (2021-10-18)
===

### Fixes
* Creating new client
* client method, OAuth 1.0 method

v0.5.0 (2021-10-17)
===

### New Supported APIs
* `GET /2/users/:id/blocking`
* `GET /2/tweets/counts/all`
* `GET /2/tweets/counts/recent`
* `GET /2/tweets/search/all API`

### Features
* Support OAuth 1.0a

### Fixes
* Resolving endpoint
* ParameterMap method
* name of some structs

### Documentation

v0.4.2 (2021-10-14)
===

### Features
* Handling rate limit error

v0.4.1 (2021-10-13)
===

### Features
* Handling not 200 errors

v0.4.0 (2021-10-12)
===

### New Supported APIs
* GET /2/tweets/search/recent

v0.3.0 (2021-10-11)
===

### New Supported APIs
* `GET /2/users/:id/followers`
* `GET /2/users/:id/following`

### Fixes
* Name of some files


v0.2.0 (2021-10-09)
====

### New Supported APIs
* `GET /2/tweets/:id`
* `GET /2/tweets`

### Features
* Support partial error
* Handle non 200 error

### Fixes
* Calling Twitter API
* User Lookup APIs
* Directory struct

### Documentation
* add CREDITS
* add LICENCE

v0.1.0 (2021-10-08)
====

* dev release 🚀