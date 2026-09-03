package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNewConfiguresHTTPServer(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New(":8080", WithLogger(logger))

	if server.httpServer.Addr != ":8080" {
		t.Fatalf("address = %q, want %q", server.httpServer.Addr, ":8080")
	}
	if server.httpServer.ReadTimeout != 10*time.Second || server.httpServer.WriteTimeout != 30*time.Second {
		t.Fatalf("unexpected server timeouts: read=%s write=%s", server.httpServer.ReadTimeout, server.httpServer.WriteTimeout)
	}
}

func TestRegisterRoutesAndUse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New(":8080", WithLogger(logger))
	server.RegisterRoutes("GET /health", http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusAccepted)
	}))
	server.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("X-Test", "true")
			next.ServeHTTP(writer, request)
		})
	})

	directRecorder := httptest.NewRecorder()
	directRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	server.ServerHttp(directRecorder, directRequest)

	if directRecorder.Code != http.StatusAccepted {
		t.Fatalf("direct status code = %d, want %d", directRecorder.Code, http.StatusAccepted)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusAccepted)
	}
	if recorder.Header().Get("X-Test") != "true" {
		t.Fatal("middleware was not applied")
	}
}

func TestStartReturnsListenError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New("invalid-address", WithLogger(logger))

	if err := server.Start(); err == nil {
		t.Fatal("Start returned nil error, want listen error")
	}
}

func TestStartShutsDownOnSignal(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New("127.0.0.1:0", WithLogger(logger))
	result := make(chan error, 1)
	go func() {
		result <- server.Start()
	}()

	time.Sleep(50 * time.Millisecond)
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("find test process: %v", err)
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("send SIGTERM: %v", err)
	}

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Start returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Start did not shut down after SIGTERM")
	}
}
