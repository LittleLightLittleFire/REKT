package fields

type SpaceField string

const (
	SpaceFieldHostIDs          SpaceField = "host_ids"
	SpaceFieldCreatedAt        SpaceField = "created_at"
	SpaceFieldCreatorID        SpaceField = "creator_id"
	SpaceFieldID               SpaceField = "id"
	SpaceFieldLang             SpaceField = "lang"
	SpaceFieldInvitedUserIDs   SpaceField = "invited_user_ids"
	SpaceFieldParticipantCount SpaceField = "participant_count"
	SpaceFieldSpeakerIDs       SpaceField = "speaker_ids"
	SpaceFieldStartedAt        SpaceField = "started_at"
	SpaceFieldState            SpaceField = "state"
	SpaceFieldTitle            SpaceField = "title"
	SpaceFieldUpdatedAt        SpaceField = "updated_at"
	SpaceFieldScheduledStart   SpaceField = "scheduled_start"
	SpaceFieldIsTicketed       SpaceField = "is_ticketed"
)

func (f SpaceField) String() string {
	return string(f)
}

type SpaceFieldList []SpaceField

func (fl SpaceFieldList) FieldsName() string {
	return "space.fields"
}

func (fl SpaceFieldList) Values() []string {
	if fl == nil {
		return []string{}
	}

	s := []string{}
	for _, f := range fl {
		s = append(s, f.String())
	}

	return s
}
