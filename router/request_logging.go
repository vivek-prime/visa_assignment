package router

import (
	"context"
	"time"

	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDKey    = utils.RequestIDKey
	RequestIDHeader = "X-Request-ID"
)

func RequestIDFromContext(ctx context.Context) (string, bool) {
	return utils.RequestIDFromContext(ctx)
}

func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Set(RequestIDKey, requestID)
		c.Request = c.Request.WithContext(utils.WithRequestID(c.Request.Context(), requestID))
		c.Header(RequestIDHeader, requestID)

		startedAt := time.Now()
		method, path := c.Request.Method, c.Request.URL.Path
		utils.LogRequest(c.Request.Context(), "request started method=%q path=%q", method, path)
		defer func() {
			utils.LogRequest(
				c.Request.Context(),
				"request completed method=%q path=%q status=%d duration=%s",
				method,
				path,
				c.Writer.Status(),
				time.Since(startedAt),
			)
		}()
		c.Next()
	}
}
