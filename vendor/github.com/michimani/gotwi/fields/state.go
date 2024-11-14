package fields

type State string

const (
	StateAll       State = "all"
	StateLive      State = "live"
	StateScheduled State = "scheduled"
)

func (s State) String() string {
	return string(s)
}

func (s State) Valid() bool {
	return s == StateAll || s == StateLive || s == StateScheduled
}
