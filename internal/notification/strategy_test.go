package notification

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestLogNotifierNotify(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	notifier := NewLogNotifier(logger)

	if err := notifier.Notify(context.Background(), "user-123", "hello"); err != nil {
		t.Fatalf("Notify() error = %v, want nil", err)
	}

	output := buf.String()
	if !strings.Contains(output, "notificación enviada") {
		t.Fatalf("log output = %q, want notification message", output)
	}
	if !strings.Contains(output, "user-123") {
		t.Fatalf("log output = %q, want user_id=user-123", output)
	}
	if !strings.Contains(output, "hello") {
		t.Fatalf("log output = %q, want message=hello", output)
	}
}

func TestNoOpNotifierNotify(t *testing.T) {
	notifier := &NoOpNotifier{}

	if err := notifier.Notify(context.Background(), "user-456", "ignored"); err != nil {
		t.Fatalf("Notify() error = %v, want nil", err)
	}
}
