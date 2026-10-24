package router

import (
	"context"
	"log"
	"time"

	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDKey    = utils.RequestIDKey
	RequestIDHeader = "X-Request-ID"
)

type requestIDContextKey struct{}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}

func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Set(RequestIDKey, requestID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), requestIDContextKey{}, requestID))
		c.Header(RequestIDHeader, requestID)

		startedAt := time.Now()
		method, path := c.Request.Method, c.Request.URL.Path
		log.Printf("request started %s=%s method=%q path=%q", RequestIDKey, requestID, method, path)
		defer func() {
			log.Printf(
				"request completed %s=%s method=%q path=%q status=%d duration=%s",
				RequestIDKey,
				requestID,
				method,
				path,
				c.Writer.Status(),
				time.Since(startedAt),
			)
		}()

		c.Next()
	}
}
