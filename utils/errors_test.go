package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mssql "github.com/microsoft/go-mssqldb"
)

func TestHandleServiceErrorMapsSQLServerErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantMessage    string
	}{
		{
			name:           "duplicate index",
			err:            mssql.Error{Number: 2601, Message: "duplicate key details"},
			wantStatusCode: http.StatusConflict,
			wantMessage:    "a record with these values already exists",
		},
		{
			name:           "duplicate constraint wrapped",
			err:            errors.Join(errors.New("repository insert failed"), mssql.Error{Number: 2627, Message: "duplicate key details"}),
			wantStatusCode: http.StatusConflict,
			wantMessage:    "a record with these values already exists",
		},
		{
			name:           "lock timeout",
			err:            mssql.Error{Number: 1222, Message: "lock timeout details"},
			wantStatusCode: http.StatusConflict,
			wantMessage:    ErrConcurrentRedemption.Error(),
		},
		{
			name:           "constraint violation",
			err:            mssql.Error{Number: 547, Message: "constraint details"},
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    "request violates a database constraint",
		},
		{
			name:           "unknown SQL Server error",
			err:            mssql.Error{Number: 1205, Message: "sensitive database details"},
			wantStatusCode: http.StatusInternalServerError,
			wantMessage:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			HandleServiceError(ctx, tt.err)

			if recorder.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d", recorder.Code, tt.wantStatusCode)
			}
			var response ErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}
			if response.Error != tt.wantMessage {
				t.Errorf("error message = %q, want %q", response.Error, tt.wantMessage)
			}
		})
	}
}

func TestHandleServiceErrorMapsDomainErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantMessage    string
	}{
		{
			name:           "member not found",
			err:            ErrMemberNotFound,
			wantStatusCode: http.StatusNotFound,
			wantMessage:    ErrMemberNotFound.Error(),
		},
		{
			name:           "invalid points",
			err:            ErrInvalidPoints,
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    ErrInvalidPoints.Error(),
		},
		{
			name:           "insufficient balance",
			err:            ErrInsufficientBalance,
			wantStatusCode: http.StatusInternalServerError,
			wantMessage:    ErrInsufficientBalance.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			HandleServiceError(ctx, tt.err)

			if recorder.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d", recorder.Code, tt.wantStatusCode)
			}
			var response ErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}
			if response.Error != tt.wantMessage {
				t.Errorf("error message = %q, want %q", response.Error, tt.wantMessage)
			}
		})
	}
}
