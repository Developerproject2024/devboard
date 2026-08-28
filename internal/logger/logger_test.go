package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestConfigs(t *testing.T) {
	if config := DefaultConfig(); config.Level != slog.LevelDebug || config.JSON || config.Output == nil {
		t.Fatalf("unexpected default config: %+v", config)
	}

	if config := ProductionConfig(); config.Level != slog.LevelInfo || !config.JSON || config.Output == nil {
		t.Fatalf("unexpected production config: %+v", config)
	}
}

func TestNewTextLogger(t *testing.T) {
	var output bytes.Buffer
	logger := New(Config{Level: slog.LevelDebug, Output: &output})

	logger.Info("request completed")

	if !strings.Contains(output.String(), "request completed") {
		t.Fatalf("log output = %q, want message", output.String())
	}
}

func TestNewJSONLoggerWithDefaultOutput(t *testing.T) {
	if logger := New(Config{JSON: true}); logger == nil {
		t.Fatal("New returned nil logger")
	}
}
