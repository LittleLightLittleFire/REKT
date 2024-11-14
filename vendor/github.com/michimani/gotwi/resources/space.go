package resources

import (
	"time"
)

type Space struct {
	ID               *string    `json:"id"`
	HostIDs          []*string  `json:"host_ids,omitempty"`
	CreatorID        *string    `json:"creator_id,omitempty"`
	Lang             *string    `json:"lang,omitempty"`
	IsTicketed       *bool      `json:"is_ticketed,omitempty"`
	InvitedUserIDs   []*string  `json:"invited_user_ids,omitempty"`
	ParticipantCount *int       `json:"participant_count,omitempty"`
	SpeakerIDs       []*string  `json:"speaker_ids,omitempty"`
	State            *string    `json:"state"`
	Title            *string    `json:"title,omitempty"`
	ScheduledStart   *time.Time `json:"scheduled_start,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	StartedAt        *time.Time `json:"started_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}
