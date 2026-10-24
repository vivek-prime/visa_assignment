package utils

import (
	"context"
	"log"

	"visa_assignment/constants"
)

const RequestIDKey = constants.RequestIDKey

type requestIDContextKey struct{}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}

func LogRequest(ctx context.Context, format string, args ...any) {
	requestID, ok := RequestIDFromContext(ctx)
	if !ok || requestID == constants.Empty {
		requestID = "NA"
	}

	log.Printf("%s=%s "+format, append([]any{RequestIDKey, requestID}, args...)...)
}
