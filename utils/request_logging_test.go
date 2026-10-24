package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestWithRequestIDAddsZerologContext(t *testing.T) {
	var output bytes.Buffer
	previousLogger := log.Logger
	log.Logger = zerolog.New(&output)
	t.Cleanup(func() {
		log.Logger = previousLogger
	})

	ctx := WithRequestID(context.Background(), "test-request-id")
	zerolog.Ctx(ctx).Info().Int("status", 200).Msg("Request completed")

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}
	if got := entry[RequestIDKey]; got != "test-request-id" {
		t.Errorf("request ID field = %v, want %q", got, "test-request-id")
	}
	if got := entry["status"]; got != float64(200) {
		t.Errorf("status = %v, want %d", got, 200)
	}
	if got := entry["message"]; got != "Request completed" {
		t.Errorf("message = %v, want %q", got, "Request completed")
	}
}
