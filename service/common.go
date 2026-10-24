package service

import "visa_assignment/constants"

// isValidPointType checks if the point_type_id is between 1 and 4
func isValidPointType(pointTypeID int) bool {
	switch pointTypeID {
	case constants.PointTypePurchaseEarning,
		constants.PointTypeReferralBonus,
		constants.PointTypeCashback,
		constants.PointTypeRedemption:
		return true
	default:
		return false
	}
}

// isRedemption returns true if the point type is a debit/redemption (point_type_id == 4)
func isRedemption(pointTypeID int) bool {
	return pointTypeID == constants.PointTypeRedemption
}
