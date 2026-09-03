package server

import (
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg.readTimeout != 10*time.Second {
		t.Fatalf("readTimeout = %s, want %s", cfg.readTimeout, 10*time.Second)
	}
	if cfg.writeTimeout != 30*time.Second {
		t.Fatalf("writeTimeout = %s, want %s", cfg.writeTimeout, 30*time.Second)
	}
	if cfg.idleTimeout != 60*time.Second {
		t.Fatalf("idleTimeout = %s, want %s", cfg.idleTimeout, 60*time.Second)
	}
	if cfg.readHeaderTimeout != 5*time.Second {
		t.Fatalf("readHeaderTimeout = %s, want %s", cfg.readHeaderTimeout, 5*time.Second)
	}
	if cfg.logger == nil {
		t.Fatal("logger = nil, want default logger")
	}
}

func TestWithOptions(t *testing.T) {
	cfg := defaultConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	WithReadTimeout(1 * time.Second)(&cfg)
	WithWriteTimeout(2 * time.Second)(&cfg)
	WithIdleTimeout(3 * time.Second)(&cfg)
	WithReadHeaderTimeout(4 * time.Second)(&cfg)
	WithLogger(logger)(&cfg)

	if cfg.readTimeout != 1*time.Second {
		t.Fatalf("readTimeout = %s, want %s", cfg.readTimeout, 1*time.Second)
	}
	if cfg.writeTimeout != 2*time.Second {
		t.Fatalf("writeTimeout = %s, want %s", cfg.writeTimeout, 2*time.Second)
	}
	if cfg.idleTimeout != 3*time.Second {
		t.Fatalf("idleTimeout = %s, want %s", cfg.idleTimeout, 3*time.Second)
	}
	if cfg.readHeaderTimeout != 4*time.Second {
		t.Fatalf("readHeaderTimeout = %s, want %s", cfg.readHeaderTimeout, 4*time.Second)
	}
	if cfg.logger != logger {
		t.Fatal("logger was not applied")
	}
}
