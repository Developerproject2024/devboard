package domain

import (
	"errors"
	"testing"
)

func TestNewTask(t *testing.T) {
	task := NewTask("project-1", "Tarea", "Descripción", "user-1")

	if task.ProjectID != "project-1" || task.Title != "Tarea" || task.Description != "Descripción" || task.CreatedBy != "user-1" {
		t.Fatalf("NewTask() devolvió una tarea inesperada: %+v", task)
	}
	if task.Status != TaskStatusTodo || task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
		t.Fatalf("NewTask() no estableció los valores iniciales: %+v", task)
	}
}

func TestTaskAssignTo(t *testing.T) {
	task := NewTask("project-1", "Tarea", "", "user-1")

	if err := task.AssignTo("user-2"); err != nil {
		t.Fatalf("AssignTo() devolvió un error inesperado: %v", err)
	}
	if task.AssigneeID == nil || *task.AssigneeID != "user-2" {
		t.Fatalf("AssignTo() no asignó el usuario esperado: %+v", task.AssigneeID)
	}
}

func TestTaskAssignToDone(t *testing.T) {
	task := NewTask("project-1", "Tarea", "", "user-1")
	task.Status = TaskStatusDone

	if err := task.AssignTo("user-2"); !errors.Is(err, ErrTaskNotAssignable) {
		t.Fatalf("AssignTo() devolvió %v; se esperaba ErrTaskNotAssignable", err)
	}
}

func TestTaskMarkAsDone(t *testing.T) {
	task := NewTask("project-1", "Tarea", "", "user-1")

	if err := task.MarkAsDone(); err != nil {
		t.Fatalf("MarkAsDone() devolvió un error inesperado: %v", err)
	}
	if task.Status != TaskStatusDone {
		t.Fatalf("MarkAsDone() dejó el estado %q", task.Status)
	}
}

func TestTaskMarkAsDoneTwice(t *testing.T) {
	task := NewTask("project-1", "Tarea", "", "user-1")
	_ = task.MarkAsDone()

	if err := task.MarkAsDone(); !errors.Is(err, ErrTaskAlreadyDone) {
		t.Fatalf("MarkAsDone() devolvió %v; se esperaba ErrTaskAlreadyDone", err)
	}
}

func TestTaskMoveToStatus(t *testing.T) {
	task := NewTask("project-1", "Tarea", "", "user-1")

	if err := task.MoveToStatus(TaskStatusInProgress); err != nil {
		t.Fatalf("MoveToStatus() devolvió un error inesperado: %v", err)
	}
	if task.Status != TaskStatusInProgress {
		t.Fatalf("MoveToStatus() dejó el estado %q", task.Status)
	}
}

func TestTaskMoveToInvalidStatus(t *testing.T) {
	task := NewTask("project-1", "Tarea", "", "user-1")

	if err := task.MoveToStatus(TaskStatus("invalid")); err == nil {
		t.Fatal("MoveToStatus() debía rechazar un estado inválido")
	}
}
