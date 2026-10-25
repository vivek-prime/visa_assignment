package service

import "visa_assignment/config"

func isPointTypeDebit(pointTypeID int, pointTypeMap map[int]config.PointTypeConfig) bool {
	return pointTypeMap[pointTypeID] == config.PointTypeDebit
}

func isValidPointType(pointTypeID int, pointTypeMap map[int]config.PointTypeConfig) bool {
	if _, ok := pointTypeMap[pointTypeID]; !ok {
		return false
	}

	return true
}
