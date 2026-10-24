package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrMemberNotFound       = errors.New("member not found")
	ErrConcurrentRedemption = errors.New("another transaction is currently in progress for this member, please retry")
	ErrInvalidPointType     = errors.New("invalid point_type_id: must be between 1 and 4")
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
		RespondWithError(c, http.StatusInternalServerError, err.Error())
	}
}
