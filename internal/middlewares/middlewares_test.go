package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryReturnsInternalServerError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("test panic")
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)

	Recovery(logger)(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(recorder.Body.String(), "error interno del servidor") {
		t.Fatalf("response body = %q", recorder.Body.String())
	}
}

func TestLoggerRecordsResponse(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte("created"))
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/items", nil)

	Logger(logger)(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if recorder.Body.String() != "created" {
		t.Fatalf("response body = %q, want %q", recorder.Body.String(), "created")
	}
	for _, expected := range []string{"request completado", "method=POST", "path=/items", "status=201"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("log output = %q, want %q", output.String(), expected)
		}
	}
}

func TestResponseRecorderWritesOnlyFirstHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	responseRecorder := newResponseRecorder(recorder)

	responseRecorder.WriteHeader(http.StatusCreated)
	responseRecorder.WriteHeader(http.StatusNoContent)

	if responseRecorder.statusCode != http.StatusCreated {
		t.Fatalf("status code = %d, want %d", responseRecorder.statusCode, http.StatusCreated)
	}
}
