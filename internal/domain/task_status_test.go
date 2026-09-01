package domain

import (
	"encoding/json"
	"testing"
)

type CreateTaskRequest struct {
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
}

func TestCreateTaskRequestWithValidStatus(t *testing.T) {
	body := []byte(`{"title": "Mi tarea", "status": "todo"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)
	if err != nil {
		t.Fatalf("no esperaba error, pero recibí: %v", err)
	}

	if req.Status != TaskStatusTodo {
		t.Fatalf("esperaba status todo, recibí: %s", req.Status)
	}
}

func TestCreateTaskRequestWithInvalidStatus(t *testing.T) {
	body := []byte(`{"title": "Mi tarea", "status": "banana"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)
	if err == nil {
		t.Fatalf("esperaba error, pero recibí nil")
	}

	expected := `estado de tarea inválido: "banana"`

	if err.Error() != expected {
		t.Fatalf("esperaba error %q, recibí %q", expected, err.Error())
	}
}

// TestTaskStatusIsValidWithValidStatuses tests that all valid statuses return true
func TestTaskStatusIsValidWithValidStatuses(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"TaskStatusTodo is valid", TaskStatusTodo, true},
		{"TaskStatusInProgress is valid", TaskStatusInProgress, true},
		{"TaskStatusDone is valid", TaskStatusDone, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.expected {
				t.Fatalf("esperaba %v, recibí %v", tt.expected, got)
			}
		})
	}
}

// TestTaskStatusIsValidWithInvalidStatuses tests that invalid statuses return false
func TestTaskStatusIsValidWithInvalidStatuses(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"Empty status is invalid", "", false},
		{"Random status is invalid", "invalid_status", false},
		{"Typo in todo", "todo_", false},
		{"Typo in in_progress", "in_progres", false},
		{"Typo in done", "done_", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.expected {
				t.Fatalf("esperaba %v, recibí %v", tt.expected, got)
			}
		})
	}
}

// TestTaskStatusUnmarshalJSONWithAllValidStatuses tests unmarshaling all valid statuses
func TestTaskStatusUnmarshalJSONWithAllValidStatuses(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected TaskStatus
	}{
		{"Unmarshal todo", `"todo"`, TaskStatusTodo},
		{"Unmarshal in_progress", `"in_progress"`, TaskStatusInProgress},
		{"Unmarshal done", `"done"`, TaskStatusDone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status TaskStatus
			err := json.Unmarshal([]byte(tt.json), &status)
			if err != nil {
				t.Fatalf("no esperaba error, pero recibí: %v", err)
			}
			if status != tt.expected {
				t.Fatalf("esperaba %q, recibí %q", tt.expected, status)
			}
		})
	}
}

// TestTaskStatusUnmarshalJSONWithInvalidStatus tests that invalid statuses return proper error
func TestTaskStatusUnmarshalJSONWithInvalidStatus(t *testing.T) {
	var status TaskStatus
	err := json.Unmarshal([]byte(`"invalid"`), &status)

	if err == nil {
		t.Fatalf("esperaba error, pero recibí nil")
	}

	expected := `estado de tarea inválido: "invalid"`
	if err.Error() != expected {
		t.Fatalf("esperaba error %q, recibí %q", expected, err.Error())
	}
}

// TestTaskStatusUnmarshalJSONWithMalformedJSON tests unmarshaling malformed JSON
func TestTaskStatusUnmarshalJSONWithMalformedJSON(t *testing.T) {
	var status TaskStatus
	err := json.Unmarshal([]byte(`not valid json`), &status)

	if err == nil {
		t.Fatalf("esperaba error, pero recibí nil")
	}
}
