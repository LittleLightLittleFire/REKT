package fields

type Exclude string

const (
	ExcludeRetweets Exclude = "retweets"
	ExcludeReplies  Exclude = "replies"
)

func (e Exclude) String() string {
	return string(e)
}

type ExcludeList []Exclude

func (el ExcludeList) FieldsName() string {
	return "exclude"
}

func (el ExcludeList) Values() []string {
	if el == nil {
		return []string{}
	}

	s := []string{}
	for _, e := range el {
		s = append(s, e.String())
	}

	return s
}
