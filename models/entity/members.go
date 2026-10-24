package model

import (
	"time"
)

// MemberEntity represents a row in the Members table and the response for POST /members
type MemberEntity struct {
	MemberID  int64     `json:"member_id" db:"member_id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
