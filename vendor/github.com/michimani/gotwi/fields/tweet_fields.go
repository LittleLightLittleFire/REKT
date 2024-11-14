package fields

type TweetField string

const (
	TweetFieldAttachments        TweetField = "attachments"
	TweetFieldAuthorID           TweetField = "author_id"
	TweetFieldContextAnnotations TweetField = "context_annotations"
	TweetFieldConversationID     TweetField = "conversation_id"
	TweetFieldCreatedAt          TweetField = "created_at"
	TweetFieldEntities           TweetField = "entities"
	TweetFieldGeo                TweetField = "geo"
	TweetFieldID                 TweetField = "id"
	TweetFieldInReplyToUserID    TweetField = "in_reply_to_user_id"
	TweetFieldLang               TweetField = "lang"
	TweetFieldNonPublicMetrics   TweetField = "non_public_metrics"
	TweetFieldPublicMetrics      TweetField = "public_metrics"
	TweetFieldOrganicMetrics     TweetField = "organic_metrics"
	TweetFieldPromotedMetrics    TweetField = "promoted_metrics"
	TweetFieldPossiblySensitive  TweetField = "possibly_sensitive"
	TweetFieldReferencedTweets   TweetField = "referenced_tweets"
	TweetFieldReplySettings      TweetField = "reply_settings"
	TweetFieldSource             TweetField = "source"
	TweetFieldText               TweetField = "text"
	TweetFieldWithheld           TweetField = "withheld"
	TweetFieldNoteTweet          TweetField = "note_tweet"
)

func (f TweetField) String() string {
	return string(f)
}

type TweetFieldList []TweetField

func (fl TweetFieldList) FieldsName() string {
	return "tweet.fields"
}

func (fl TweetFieldList) Values() []string {
	if fl == nil {
		return []string{}
	}

	s := []string{}
	for _, f := range fl {
		s = append(s, f.String())
	}

	return s
}
