package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Developerproject2024/devboard/internal/domain"
	"github.com/Developerproject2024/devboard/internal/repository/memory"
	"github.com/Developerproject2024/devboard/internal/usecase"
	"github.com/Developerproject2024/devboard/internal/validator"
)

func newTaskHandlerForTest() *TaskHandler {
	repository := memory.NewTaskRepository()
	return NewTaskHandler(
		usecase.NewTaskUseCase(repository, repository),
		validator.New(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
}

func taskRequest(method, path, body string) *http.Request {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.SetPathValue("id", "task_1")
	return request
}

func TestNewTaskHandler(t *testing.T) {
	handler := newTaskHandlerForTest()
	if handler == nil || handler.validator == nil || handler.logger == nil {
		t.Fatal("NewTaskHandler() no inicializó el handler")
	}
}

func TestTaskHandlerCreate(t *testing.T) {
	handler := newTaskHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Create(recorder, taskRequest(http.MethodPost, "/tasks", `{"project_id":"project-1","title":"Tarea","description":"Descripción"}`))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("Create() status = %d; se esperaba %d", recorder.Code, http.StatusCreated)
	}
	var response taskResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Create() devolvió JSON inválido: %v", err)
	}
	if response.ID == "" || response.Title != "Tarea" {
		t.Fatalf("Create() devolvió una respuesta inesperada: %+v", response)
	}
}

func TestTaskHandlerCreateInvalidRequest(t *testing.T) {
	handler := newTaskHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Create(recorder, taskRequest(http.MethodPost, "/tasks", `{`))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Create() status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestTaskHandlerGet(t *testing.T) {
	handler := newTaskHandlerForTest()
	createRecorder := httptest.NewRecorder()
	handler.Create(createRecorder, taskRequest(http.MethodPost, "/tasks", `{"project_id":"project-1","title":"Tarea"}`))

	recorder := httptest.NewRecorder()
	handler.Get(recorder, taskRequest(http.MethodGet, "/tasks/task_1", ""))

	if recorder.Code != http.StatusOK {
		t.Fatalf("Get() status = %d; se esperaba %d", recorder.Code, http.StatusOK)
	}
}

func TestTaskHandlerGetNotFound(t *testing.T) {
	handler := newTaskHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Get(recorder, taskRequest(http.MethodGet, "/tasks/task_1", ""))

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("Get() status = %d; se esperaba %d", recorder.Code, http.StatusNotFound)
	}
}

func TestTaskHandlerListByProject(t *testing.T) {
	handler := newTaskHandlerForTest()
	handler.Create(httptest.NewRecorder(), taskRequest(http.MethodPost, "/tasks", `{"project_id":"project-1","title":"Tarea"}`))

	recorder := httptest.NewRecorder()
	request := taskRequest(http.MethodGet, "/projects/project-1/tasks", "")
	request.SetPathValue("id", "project-1")
	handler.ListByProject(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("ListByProject() status = %d; se esperaba %d", recorder.Code, http.StatusOK)
	}
	var response paginatedTaskResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("ListByProject() devolvió JSON inválido: %v", err)
	}
	if len(response.Data) != 1 {
		t.Fatalf("ListByProject() devolvió %d tareas; se esperaba 1", len(response.Data))
	}
}

func TestTaskHandlerAssign(t *testing.T) {
	handler := newTaskHandlerForTest()
	handler.Create(httptest.NewRecorder(), taskRequest(http.MethodPost, "/tasks", `{"project_id":"project-1","title":"Tarea"}`))

	recorder := httptest.NewRecorder()
	handler.Assign(recorder, taskRequest(http.MethodPut, "/tasks/task_1/assign", `{"assignee_id":"550e8400-e29b-41d4-a716-446655440000"}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("Assign() status = %d; se esperaba %d", recorder.Code, http.StatusOK)
	}
}

func TestTaskHandlerAssignInvalidRequest(t *testing.T) {
	handler := newTaskHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Assign(recorder, taskRequest(http.MethodPut, "/tasks/task_1/assign", `{}`))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Assign() status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestTaskHandlerUpdateStatus(t *testing.T) {
	handler := newTaskHandlerForTest()
	handler.Create(httptest.NewRecorder(), taskRequest(http.MethodPost, "/tasks", `{"project_id":"project-1","title":"Tarea"}`))

	recorder := httptest.NewRecorder()
	handler.UpdateStatus(recorder, taskRequest(http.MethodPut, "/tasks/task_1/status", `{"status":"in_progress"}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("UpdateStatus() status = %d; se esperaba %d", recorder.Code, http.StatusOK)
	}
	var response taskResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("UpdateStatus() devolvió JSON inválido: %v", err)
	}
	if response.Status != domain.TaskStatusInProgress {
		t.Fatalf("UpdateStatus() dejó estado %q", response.Status)
	}
}

func TestTaskHandlerUpdateStatusInvalidRequest(t *testing.T) {
	handler := newTaskHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.UpdateStatus(recorder, taskRequest(http.MethodPut, "/tasks/task_1/status", `{}`))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("UpdateStatus() status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}
