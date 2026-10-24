package router

import (
	"context"
	"time"

	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
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
		zerolog.Ctx(c.Request.Context()).Info().
			Str("method", method).
			Str("path", path).
			Msg("Request started")
		defer func() {
			zerolog.Ctx(c.Request.Context()).Info().
				Str("method", method).
				Str("path", path).
				Int("status", c.Writer.Status()).
				Dur("duration", time.Since(startedAt)).
				Msg("Request completed")
		}()
		c.Next()
	}
}
