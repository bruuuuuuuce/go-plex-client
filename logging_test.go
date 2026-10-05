package plex

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlexStructuredLogger(t *testing.T) {
	var output bytes.Buffer
	client := Plex{}
	client.SetLogger(slog.New(slog.NewJSONHandler(&output, nil)))

	client.logger().Error("request failed", "error", errors.New("boom"))

	assert.Contains(t, output.String(), `"level":"ERROR"`)
	assert.Contains(t, output.String(), `"msg":"request failed"`)
	assert.Contains(t, output.String(), `"error":"boom"`)
}

func TestLoggingDisabledByDefault(t *testing.T) {
	client := Plex{}

	assert.NotPanics(t, func() {
		client.logger().Info("discarded")
	})
}

func TestWebhookStructuredLogger(t *testing.T) {
	var output bytes.Buffer
	webhooks := NewWebhook()
	webhooks.SetLogger(slog.New(slog.NewJSONHandler(&output, nil)))

	webhooks.logger().Warn("unknown event", "event", "example")

	assert.Contains(t, output.String(), `"level":"WARN"`)
	assert.Contains(t, output.String(), `"event":"example"`)
}
