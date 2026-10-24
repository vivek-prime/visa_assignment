package model

import "time"

// MemberRequest represents the request body for POST /members
type MemberRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// MemberBalanceResponse represents the response body for GET /members/:memberId
type MemberBalanceResponse struct {
	MemberID      int64     `json:"member_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PointsBalance int64     `json:"points_balance"`
	CreatedAt     time.Time `json:"created_at"`
}
