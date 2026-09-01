package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Developerproject2024/devboard/internal/domain"
	"github.com/Developerproject2024/devboard/internal/validator"
)

func TestResponseJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	err := ResponseJSON(recorder, http.StatusOK, data)

	if err != nil {
		t.Fatalf("ResponseJSON() error = %v, want nil", err)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("message = %q, want %q", response["message"], "success")
	}
}

func TestResponseJSONWithNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := ResponseJSON(recorder, http.StatusNotFound, map[string]string{"error": "not found"})

	if err != nil {
		t.Fatalf("ResponseJSON() error = %v, want nil", err)
	}

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestRespondValidationError(t *testing.T) {
	recorder := httptest.NewRecorder()
	validationErrors := []validator.ValidationError{
		{Field: "email", Message: "invalid email"},
		{Field: "name", Message: "name is required"},
	}

	err := RespondValidationError(recorder, validationErrors)

	if err != nil {
		t.Fatalf("RespondValidationError() error = %v, want nil", err)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusBadRequest)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	var response errorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Error != "datos de entrada inválidos" {
		t.Errorf("error = %q, want %q", response.Error, "datos de entrada inválidos")
	}

	if len(response.Details) != 2 {
		t.Errorf("details length = %d, want %d", len(response.Details), 2)
	}
}

func TestRespondErrorNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	logger := slog.Default()

	RespondError(recorder, request, logger, domain.ErrNotFound)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/problem+json")
	}

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Status != http.StatusNotFound {
		t.Errorf("problem.Status = %d, want %d", problem.Status, http.StatusNotFound)
	}

	if problem.Title != "Recurso no encontrado" {
		t.Errorf("problem.Title = %q, want %q", problem.Title, "Recurso no encontrado")
	}
}

func TestRespondErrorAlreadyExists(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	logger := slog.Default()

	RespondError(recorder, request, logger, domain.ErrAlreadyExists)

	if recorder.Code != http.StatusConflict {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusConflict)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Status != http.StatusConflict {
		t.Errorf("problem.Status = %d, want %d", problem.Status, http.StatusConflict)
	}

	if problem.Title != "El recurso ya existe" {
		t.Errorf("problem.Title = %q, want %q", problem.Title, "El recurso ya existe")
	}
}

func TestRespondErrorUnauthorized(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	logger := slog.Default()

	RespondError(recorder, request, logger, domain.ErrUnauthorized)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Title != "No autorizado" {
		t.Errorf("problem.Title = %q, want %q", problem.Title, "No autorizado")
	}
}

func TestRespondErrorForbidden(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	logger := slog.Default()

	RespondError(recorder, request, logger, domain.ErrForbidden)

	if recorder.Code != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusForbidden)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Title != "Recurso prohíbido" {
		t.Errorf("problem.Title = %q, want %q", problem.Title, "Recurso prohíbido")
	}
}

func TestRespondErrorInvalidInput(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	logger := slog.Default()

	RespondError(recorder, request, logger, domain.ErrInvalidInput)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusBadRequest)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Title != "Solicitud no válida" {
		t.Errorf("problem.Title = %q, want %q", problem.Title, "Solicitud no válida")
	}
}

func TestRespondErrorUnknown(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	logger := slog.Default()

	unknownErr := errors.New("unexpected error")
	RespondError(recorder, request, logger, unknownErr)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Status != http.StatusInternalServerError {
		t.Errorf("problem.Status = %d, want %d", problem.Status, http.StatusInternalServerError)
	}

	if problem.Title != "Error interno del servidor" {
		t.Errorf("problem.Title = %q, want %q", problem.Title, "Error interno del servidor")
	}
}

func TestProblemDetailsIncludesInstance(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/tasks/123", nil)
	logger := slog.Default()

	RespondError(recorder, request, logger, domain.ErrNotFound)

	var problem ProblemDetails
	if err := json.NewDecoder(recorder.Body).Decode(&problem); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if problem.Instance != "/api/tasks/123" {
		t.Errorf("problem.Instance = %q, want %q", problem.Instance, "/api/tasks/123")
	}
}
