package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandlerServeHTTP(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	NewHealthHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("content type = %q, want %q", contentType, "application/json")
	}

	var response healthResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response != (healthResponse{Status: "Fabio Arango", Version: "1.1.0"}) {
		t.Fatalf("response = %+v, want %+v", response, healthResponse{Status: "Fabio Arango", Version: "1.1.0"})
	}
}

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()
	if handler == nil {
		t.Fatalf("NewHealthHandler() should not return nil")
	}
}

func TestParsePaginationWithDefaultValues(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	params := ParsePagination(request)

	if params.Cursor != "" {
		t.Errorf("cursor = %q, want %q", params.Cursor, "")
	}
	if params.Limit != 20 {
		t.Errorf("limit = %d, want %d", params.Limit, 20)
	}
}

func TestParsePaginationWithCustomLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?limit=50", nil)
	params := ParsePagination(request)

	if params.Limit != 50 {
		t.Errorf("limit = %d, want %d", params.Limit, 50)
	}
}

func TestParsePaginationWithMaxLimitExceeded(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?limit=200", nil)
	params := ParsePagination(request)

	if params.Limit != 100 {
		t.Errorf("limit = %d, want %d (should be capped at 100)", params.Limit, 100)
	}
}

func TestParsePaginationWithInvalidLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?limit=invalid", nil)
	params := ParsePagination(request)

	if params.Limit != 20 {
		t.Errorf("limit = %d, want %d (should use default for invalid)", params.Limit, 20)
	}
}

func TestParsePaginationWithNegativeLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?limit=-10", nil)
	params := ParsePagination(request)

	if params.Limit != 20 {
		t.Errorf("limit = %d, want %d (should use default for negative)", params.Limit, 20)
	}
}

func TestParsePaginationWithCursor(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?cursor=abc123", nil)
	params := ParsePagination(request)

	if params.Cursor != "abc123" {
		t.Errorf("cursor = %q, want %q", params.Cursor, "abc123")
	}
}

func TestParsePaginationWithBothParams(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?cursor=xyz789&limit=30", nil)
	params := ParsePagination(request)

	if params.Cursor != "xyz789" {
		t.Errorf("cursor = %q, want %q", params.Cursor, "xyz789")
	}
	if params.Limit != 30 {
		t.Errorf("limit = %d, want %d", params.Limit, 30)
	}
}
