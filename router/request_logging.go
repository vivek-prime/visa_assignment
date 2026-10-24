package router

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDKey    = "request_id"
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
		log.Printf("request started request_id=%s method=%q path=%q", requestID, method, path)
		defer func() {
			log.Printf(
				"request completed request_id=%s method=%q path=%q status=%d duration=%s",
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
