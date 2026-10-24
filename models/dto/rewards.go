package model

// RewardRequest represents the request body for POST /rewards
type RewardRequest struct {
	MemberID    int64  `json:"member_id" binding:"required"`
	PointTypeID int    `json:"point_type_id" binding:"required"`
	Points      int64  `json:"points" binding:"required"`
	Description string `json:"description" binding:"required"`
}
