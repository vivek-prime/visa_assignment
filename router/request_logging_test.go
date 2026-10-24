package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggingPropagatesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLogging())
	r.GET("/test", func(c *gin.Context) {
		requestID := c.GetString(RequestIDKey)
		contextRequestID, ok := utils.RequestIDFromContext(c.Request.Context())
		if !ok || requestID == "" || requestID != contextRequestID {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	response := httptest.NewRecorder()
	r.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/test", nil))

	if response.Code != http.StatusNoContent {
		t.Fatalf("request returned status %d, want %d", response.Code, http.StatusNoContent)
	}
	if response.Header().Get(RequestIDHeader) == "" {
		t.Fatal("response is missing request ID header")
	}
}
