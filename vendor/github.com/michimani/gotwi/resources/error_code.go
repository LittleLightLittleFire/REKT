package resources

import "fmt"

type ErrorCode int

func (e ErrorCode) Detail() *ErrorCodeDetail {
	if d, ok := errorCodeMap[e]; ok {
		return &d
	}

	return &ErrorCodeDetail{
		Text:        "Undifined code.",
		Description: fmt.Sprintf("'%d' is undefined error code.", int(e)),
	}
}

type ErrorCodeDetail struct {
	Text        string
	Description string
}

var errorCodeMap map[ErrorCode]ErrorCodeDetail = map[ErrorCode]ErrorCodeDetail{
	3:   {Text: "Invalid coordinates.", Description: "Corresponds with HTTP 400. The coordinates provided as parameters were not valid for the request."},
	13:  {Text: "No location associated with the specified IP address.", Description: "Corresponds with HTTP 404. It was not possible to derive a location for the IP address provided as a parameter on the geo search request."},
	17:  {Text: "No user matches for specified terms.", Description: "Corresponds with HTTP 404. It was not possible to find a user profile matching the parameters specified."},
	32:  {Text: "Could not authenticate you.", Description: "Corresponds with HTTP 401. There was an issue with the authentication data for the request."},
	34:  {Text: "Sorry, that page does not exist.", Description: "Corresponds with HTTP 404. The specified resource was not found."},
	36:  {Text: "You cannot report yourself for spam.", Description: "Corresponds with HTTP 403. You cannot use your own user ID in a report spam call."},
	38:  {Text: "<named> parameter is missing.", Description: "Corresponds with HTTP 403. The request is missing the <named> parameter (such as media, text, etc.) in the request."},
	44:  {Text: "attachment_url parameter is invalid.", Description: "Corresponds with HTTP 400. The URL value provided is not a URL that can be attached to this Tweet."},
	50:  {Text: "User not found.", Description: "Corresponds with HTTP 404. The user is not found."},
	63:  {Text: "User has been suspended.", Description: "Corresponds with HTTP 403 The user account has been suspended and information cannot be retrieved."},
	64:  {Text: "Your account is suspended and is not permitted to access this feature.", Description: "Corresponds with HTTP 403. The access token being used belongs to a suspended user."},
	68:  {Text: "Some actions on this user's Tweet have been disabled by Twitter. The Twitter REST API v1 is no longer active. ", Description: "Corresponds with HTTP 410. The request was made to a retired v1-era URL."},
	87:  {Text: "Client is not permitted to perform this action.", Description: "Corresponds with HTTP 403. The endpoint called is not a permitted URL."},
	88:  {Text: "Rate limit exceeded.", Description: "Corresponds with HTTP 429. The request limit for this resource has been reached for the current rate limit window."},
	89:  {Text: "Invalid or expired token.", Description: "Corresponds with HTTP 403. The access token used in the request is incorrect or has expired."},
	92:  {Text: "SSL is required.", Description: "Corresponds with HTTP 403. Only TLS v1.2 connections are allowed in the API. Update the request to a secure connection. See how to connect using TLS."},
	93:  {Text: "This App is not allowed to access or delete your Direct Messages.", Description: "Corresponds with HTTP 403. The OAuth token does not provide access to Direct Messages."},
	99:  {Text: "Unable to verify your credentials.", Description: "Corresponds with HTTP 403. The OAuth credentials cannot be validated. Check that the token is still valid."},
	109: {Text: "The specified user is not found in this list.", Description: "Corresponds with HTTP 404. Not Found."},
	110: {Text: "The user you are trying to remove from this list is not a member.", Description: "Corresponds with HTTP 400. Bad Request."},
	120: {Text: "Account update failed: value is too long (maximum is nn characters).", Description: "Corresponds with HTTP 403. Thrown when one of the values passed to the update_profile.json endpoint exceeds the maximum value currently permitted for that field. The error message will specify the allowable maximum number of nn characters."},
	130: {Text: "Over capacity.", Description: "Corresponds with HTTP 503. Twitter is temporarily over capacity."},
	131: {Text: "Internal error.", Description: "Corresponds with HTTP 500. An unknown internal error occurred."},
	135: {Text: "Could not authenticate you.", Description: "Corresponds with HTTP 401. Timestamp out of bounds (often caused by a clock drift when authenticating - check your system clock)."},
	139: {Text: "You have already favorited this status.", Description: "Corresponds with HTTP 403. A Tweet cannot be favorited (liked) more than once."},
	144: {Text: "No status found with that ID.", Description: "Corresponds with HTTP 404. The requested Tweet ID is not found (if it existed, it was probably deleted)."},
	150: {Text: "You cannot send messages to users who are not following you.", Description: "Corresponds with HTTP 403. Sending a Direct Message failed."},
	151: {Text: "There was an error sending your message: reason.", Description: "Corresponds with HTTP 403. Sending a Direct Message failed. The reason value will provide more information."},
	160: {Text: "You've already requested to follow the user.", Description: "Corresponds with HTTP 403. This was a duplicated follow request and a previous request was not yet acknowleged."},
	161: {Text: "You are unable to follow more people at this time.", Description: "Corresponds with HTTP 403. Thrown when a user cannot follow another user due to reaching the limit. This limit is applied to each user individually, independent of the Apps they use to access the Twitter platform."},
	179: {Text: "Sorry, you are not authorized to see this status.", Description: "Corresponds with HTTP 403. Thrown when a Tweet cannot be viewed by the authenticating user, usually due to the Tweet’s author having protected their Tweets."},
	185: {Text: "User is over daily status update limit.", Description: "Corresponds with HTTP 403. Thrown when a Tweet cannot be posted due to the user having no allowance remaining to post. Despite the text in the error message indicating that this error is only thrown when a daily limit is reached, this error will be thrown whenever a posting limitation has been reached. Posting allowances have roaming windows of time of unspecified duration."},
	186: {Text: "Tweet needs to be a bit shorter.", Description: "Corresponds with HTTP 403. The status text is too long."},
	187: {Text: "Status is a duplicate.", Description: "Corresponds with HTTP 403. The status text has already been Tweeted by the authenticated account."},
	195: {Text: "Missing or invalid url parameter.", Description: "Corresponds with HTTP 403.  The request needs to have a valid url parameter."},
	205: {Text: "You are over the limit for spam reports.", Description: "Corresponds with HTTP 403. The account limit for reporting spam has been reached. Try again later."},
	214: {Text: "Owner must allow dms from anyone.", Description: "Corresponds with HTTP 403. The user is not set up to have open Direct Messages when trying to set up a welcome message."},
	215: {Text: "Bad authentication data.", Description: "Corresponds with HTTP 400. The method requires authentication but it was not presented or was wholly invalid."},
	220: {Text: "Your credentials do not allow access to this resource.", Description: "Corresponds with HTTP 403. The authentication token in use is restricted and cannot access the requested resource."},
	226: {Text: "This request looks like it might be automated. To protect our users from spam and other malicious activity, we can’t complete this action right now.", Description: "Corresponds with HTTP 403. We constantly monitor and adjust our filters to block spam and malicious activity on the Twitter platform. These systems are tuned in real-time. If you get this response our systems have flagged the Tweet or Direct Message as possibly fitting this profile. If you believe that the Tweet or DM you attempted to create was flagged in error, report the details by filing a ticket at https://help.twitter.com/forms/platform."},
	251: {Text: "This endpoint has been retired and should not be used.", Description: "Corresponds with HTTP 410. The App made a request to a retired URL."},
	261: {Text: "Your App cannot perform write actions.", Description: "Corresponds with HTTP 403. Caused by the App being restricted from POST, PUT, or DELETE actions. Check the information on your App dashboard. You may also file a ticket at https://help.twitter.com/forms/platform."},
	271: {Text: "You can’t mute yourself.", Description: "Corresponds with HTTP 403. The authenticated user account cannot mute itself."},
	272: {Text: "You are not muting the specified user.", Description: "Corresponds with HTTP 403. The authenticated user account is not muting the account a call is attempting to unmute."},
	323: {Text: "Animated GIFs are not allowed when uploading multiple images.", Description: "Corresponds with HTTP 400. Only one animated GIF may be attached to a single Tweet."},
	324: {Text: "The validation of media ids failed.", Description: "Corresponds with HTTP 400. There was a problem with the media ID submitted with the Tweet."},
	325: {Text: "A media id was not found.", Description: "Corresponds with HTTP 400. The media ID attached to the Tweet was not found."},
	326: {Text: "To protect our users from spam and other malicious activity, this account is temporarily locked.", Description: "Corresponds with HTTP 403. The user should log in to https://twitter.com to unlock their account before the user token can be used."},
	327: {Text: "You have already retweeted this Tweet.", Description: "Corresponds with HTTP 403. The user cannot retweet the same Tweet more than once."},
	349: {Text: "You cannot send messages to this user.", Description: "Corresponds with HTTP 403. The sender does not have privileges to Direct Message the recipient."},
	354: {Text: "The text of your direct message is over the max character limit.", Description: "Corresponds with HTTP 403. The message size exceeds the number of characters permitted in a Direct Message."},
	355: {Text: "Subscription already exists.", Description: "Corresponds with HTTP 409 Conflict. Related to Account Activity API request to add a new subscription for an authenticated user."},
	385: {Text: "You attempted to reply to a Tweet that is deleted or not visible to you.", Description: "Corresponds with HTTP 403. A reply can only be sent with reference to an existing public Tweet."},
	386: {Text: "The Tweet exceeds the number of allowed attachment types.", Description: "Corresponds with HTTP 403. A Tweet is limited to a single attachment resource (media, Quote Tweet, etc.)"},
	407: {Text: "The given URL is invalid.", Description: "Corresponds with HTTP 400. A URL included in the Tweet could not be handled. This may be because a non-ASCII URL could not be converted, or for other reasons."},
	415: {Text: "Callback URL not approved for this client App. Approved callback URLs can be adjusted in your App's settings.", Description: "Corresponds with HTTP 403. The App callback URLs must be allowlisted via the App details page in the developer portal. Only approved callback URLs may be used by the Twitter App. See the Callback URL documentation."},
	416: {Text: "Invalid / suspended App.", Description: "Corresponds with HTTP 401. The App has been suspended and cannot be used with Sign-in with Twitter."},
	417: {Text: "Desktop applications only support the oauth_callback value 'oob'.", Description: "Corresponds with HTTP 401. The App is attempting to use out-of-band PIN-based OAuth, but a callback URL has been specified in the App's settings."},
	421: {Text: "This Tweet is no longer available.", Description: "Corresponds with HTTP 404. The Tweet cannot be retrieved. This may be for a number of reasons. Read about the Twitter Rules."},
	422: {Text: "This Tweet is no longer available because it violated the Twitter Rules.", Description: "Corresponds with HTTP 404. The Tweet is not available in the API. Read about the Twitter Rules."},
	425: {Text: "Some actions on this user's Tweet have been disabled by Twitter.", Description: "Corresponds with HTTP 403. Forbidden. Read about public-interest exceptions on Twitter."},
	433: {Text: "The original Tweet author restricted who can reply to this Tweet.", Description: "Corresponds with HTTP 403. Thrown when reply to a Tweet, and the author of that original Tweet limited who can reply. In this case, a reply can only be sent if the author follows or has been mentioned by the author of the original Tweet."},
}
