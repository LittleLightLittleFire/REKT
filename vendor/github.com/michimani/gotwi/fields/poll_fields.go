package fields

type PollField string

const (
	PollFieldDurationMinutes PollField = "duration_minutes"
	PollFieldEndDatetime     PollField = "end_datetime"
	PollFieldID              PollField = "id"
	PollFieldOptions         PollField = "options"
	PollFieldVotingStatus    PollField = "voting_status"
)

func (f PollField) String() string {
	return string(f)
}

type PollFieldList []PollField

func (fl PollFieldList) FieldsName() string {
	return "poll.fields"
}

func (fl PollFieldList) Values() []string {
	if fl == nil {
		return []string{}
	}

	s := []string{}
	for _, f := range fl {
		s = append(s, f.String())
	}

	return s
}
