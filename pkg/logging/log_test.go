package logging

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Logger(t *testing.T) {
	t.Setenv("ROUTER_LOG", "true")
	l := Get()
	assert.NotNil(t, l)
	w := bytes.NewBuffer([]byte{})

	handler := slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
	slogLogger := slog.New(handler)
	Set(slogLogger)

	logger := Get()
	assert.NotNil(t, logger)
	assert.NotSame(t, l, logger)

	logger.Info("information", slog.Int("test", 1))
	logger.Debug("debug", slog.Int("test", 2))
	logger.Warn("warning", slog.Int("test", 3))

	content := w.String()

	assert.Contains(t, content, "information")
	assert.Contains(t, content, "debug")
	assert.Contains(t, content, "warning")

}
