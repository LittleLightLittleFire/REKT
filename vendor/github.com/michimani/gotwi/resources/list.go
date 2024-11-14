package resources

import "time"

type List struct {
	ID            *string    `json:"id"`
	Name          *string    `json:"name"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Private       *bool      `json:"private,omitempty"`
	FollowerCount *int       `json:"follower_count,omitempty"`
	MemberCount   *int       `json:"member_count,omitempty"`
	OwnerID       *string    `json:"owner_id,omitempty"`
	Description   *string    `json:"description,omitempty"`
}
