package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/Developerproject2024/devboard/internal/domain"
)

func TestTaskRepository_Create(t *testing.T) {
	repo := NewTaskRepository()
	task := domain.NewTask("project-1", "Implementar endpoint", "Crear el endpoint de tareas", "user-1")

	if err := repo.Create(context.Background(), task); err != nil {
		t.Fatalf("Create() devolvió un error inesperado: %v", err)
	}

	if task.ID != "task_1" {
		t.Fatalf("Create() estableció un ID %q; se esperaba %q", task.ID, "task_1")
	}

	found, err := repo.GetByID(context.Background(), task.ID)
	if err != nil {
		t.Fatalf("GetByID() devolvió un error inesperado: %v", err)
	}
	if found != task {
		t.Fatal("GetByID() no devolvió la tarea creada")
	}
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	repo := NewTaskRepository()

	_, err := repo.GetByID(context.Background(), "task_999")
	if err == nil {
		t.Fatal("GetByID() debía devolver un error")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("GetByID() devolvió %v; se esperaba ErrNotFound", err)
	}
}

func TestTaskRepository_Update(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()
	task := domain.NewTask("project-1", "Título original", "Descripción original", "user-1")

	if err := repo.Create(ctx, task); err != nil {
		t.Fatalf("no se pudo preparar la tarea: %v", err)
	}

	task.Title = "Título actualizado"
	if err := repo.Update(ctx, task); err != nil {
		t.Fatalf("Update() devolvió un error inesperado: %v", err)
	}

	found, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetByID() devolvió un error inesperado: %v", err)
	}
	if found.Title != "Título actualizado" {
		t.Errorf("Update() no persistió el título; se obtuvo %q", found.Title)
	}
}

func TestTaskRepository_ListByProject(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()
	tasks := []*domain.Task{
		domain.NewTask("project-1", "Tarea 1", "", "user-1"),
		domain.NewTask("project-1", "Tarea 2", "", "user-1"),
		domain.NewTask("project-2", "Tarea 3", "", "user-1"),
	}

	for _, task := range tasks {
		if err := repo.Create(ctx, task); err != nil {
			t.Fatalf("no se pudo preparar la tarea: %v", err)
		}
	}

	found, err := repo.ListByProject(ctx, "project-1")
	if err != nil {
		t.Fatalf("ListByProject() devolvió un error inesperado: %v", err)
	}
	if len(found) != 2 {
		t.Fatalf("ListByProject() devolvió %d tareas; se esperaban 2", len(found))
	}

	for _, task := range found {
		if task.ProjectID != "project-1" {
			t.Errorf("ListByProject() devolvió una tarea del proyecto %q", task.ProjectID)
		}
	}
}

func TestTaskRepository_ListByProject_Empty(t *testing.T) {
	repo := NewTaskRepository()

	found, err := repo.ListByProject(context.Background(), "project-unknown")
	if err != nil {
		t.Fatalf("ListByProject() devolvió un error inesperado: %v", err)
	}
	if len(found) != 0 {
		t.Fatalf("ListByProject() devolvió %d tareas; se esperaban 0", len(found))
	}
}
