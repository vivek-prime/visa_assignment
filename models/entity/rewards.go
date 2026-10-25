package model

import "time"

// RewardEntity represents a row in the rewards table
type RewardEntity struct {
	RewardID    int64     `json:"reward_id" db:"reward_id"`
	MemberID    int64     `json:"member_id" db:"member_id"`
	PointTypeID int       `json:"point_type_id" db:"point_type_id"`
	Points      int64     `json:"points" db:"points"`
	Description string    `json:"description" db:"description"`
	EventDate   time.Time `json:"event_date" db:"event_date"`
}
