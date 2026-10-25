package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	sqlServerDuplicateIndexError      int32 = 2601
	sqlServerDuplicateConstraintError int32 = 2627
	sqlServerConstraintViolationError int32 = 547
	sqlServerLockTimeoutError         int32 = 1222
)

var (
	ErrMemberNotFound       = errors.New("member not found")
	ErrDuplicateRecord      = errors.New("duplicate record found")
	ErrConcurrentRedemption = errors.New("another transaction is currently in progress for this member, please retry")
	ErrInvalidPointType     = errors.New("invalid point_type_id: no direction is configured")
	ErrInvalidPoints        = errors.New("points must be a positive number greater than zero")
	ErrInsufficientBalance  = errors.New("insufficient points balance for this redemption")
	ErrInvalidInput         = errors.New("invalid input provided")
)

// ErrorResponse represents a standardized JSON error payload
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends a JSON error response with the specified HTTP status code
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{Error: message})
}

// HandleServiceError maps domain/service errors to appropriate HTTP status codes
func HandleServiceError(c *gin.Context, err error) {
	var sqlErr interface{ SQLErrorNumber() int32 }
	if errors.As(err, &sqlErr) {
		switch sqlErr.SQLErrorNumber() {
		case sqlServerDuplicateConstraintError:
			RespondWithError(c, http.StatusInternalServerError, ErrDuplicateRecord.Error())
			return
		case sqlServerLockTimeoutError:
			RespondWithError(c, http.StatusInternalServerError, ErrConcurrentRedemption.Error())
			return
		case sqlServerConstraintViolationError:
			RespondWithError(c, http.StatusInternalServerError, ErrMemberNotFound.Error())
			return
		}
	}

	switch {
	case errors.Is(err, ErrMemberNotFound):
		RespondWithError(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrConcurrentRedemption):
		RespondWithError(c, http.StatusConflict, err.Error())
	case errors.Is(err, ErrInsufficientBalance):
		RespondWithError(c, http.StatusInternalServerError, err.Error())
	case errors.Is(err, ErrInvalidPointType),
		errors.Is(err, ErrInvalidPoints),
		errors.Is(err, ErrInvalidInput):
		RespondWithError(c, http.StatusBadRequest, err.Error())
	default:
		RespondWithError(c, http.StatusInternalServerError, "internal server error")
	}
}
