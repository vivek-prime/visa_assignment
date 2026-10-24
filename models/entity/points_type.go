package model

const (
	PointTypePurchaseEarning = iota
	PointTypeReferralBonus
	PointTypeCashback
	PointTypeRedemption
)

// PointType represents a row in the PointTypes table
type PointType struct {
	PointTypeID int    `json:"point_type_id" db:"point_type_id"`
	Description string `json:"description" db:"description"`
}

// IsValidPointType checks if the point_type_id is between 1 and 4
func IsValidPointType(pointTypeID int) bool {
	switch pointTypeID {
	case PointTypePurchaseEarning,
		PointTypeReferralBonus,
		PointTypeCashback,
		PointTypeRedemption:
		return true
	default:
		return false
	}
}

// IsRedemption returns true if the point type is a debit/redemption (point_type_id == 4)
func IsRedemption(pointTypeID int) bool {
	return pointTypeID == PointTypeRedemption
}
