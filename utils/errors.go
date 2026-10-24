package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrDuplicateEmail is returned when creating a member with an existing email
	ErrDuplicateEmail = errors.New("a member with this email already exists")

	// ErrMemberNotFound is returned when member_id does not reference an existing member
	ErrMemberNotFound = errors.New("member not found")

	ErrConcurrentRedemption = errors.New("another transaction is currently in progress for this member, please retry")

	// ErrInvalidPointType is returned when point_type_id is not between 1 and 4[cite: 4]
	ErrInvalidPointType = errors.New("invalid point_type_id: must be between 1 and 4")

	// ErrInvalidPoints is returned when points is not a positive number (> 0)[cite: 4]
	ErrInvalidPoints = errors.New("points must be a positive number greater than zero")

	// ErrInsufficientBalance is returned when a redemption exceeds the member's available balance[cite: 4]
	ErrInsufficientBalance = errors.New("insufficient points balance for this redemption")

	// ErrInvalidInput is returned when required fields are empty or malformed
	ErrInvalidInput = errors.New("invalid input provided")
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
	switch {
	case errors.Is(err, ErrDuplicateEmail):
		RespondWithError(c, http.StatusConflict, err.Error())
	case errors.Is(err, ErrMemberNotFound):
		RespondWithError(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInsufficientBalance):
		RespondWithError(c, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, ErrInvalidPointType),
		errors.Is(err, ErrInvalidPoints),
		errors.Is(err, ErrInvalidInput):
		RespondWithError(c, http.StatusBadRequest, err.Error())
	default:
		RespondWithError(c, http.StatusInternalServerError, "internal server error")
	}
}
