package utils

import (
	"context"

	"visa_assignment/constants"

	"github.com/rs/zerolog/log"
)

const RequestIDKey = constants.RequestIDKey

type requestIDContextKey struct{}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, requestIDContextKey{}, requestID)
	logger := log.Logger.With().Str(RequestIDKey, requestID).Logger()
	return logger.WithContext(ctx)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}
