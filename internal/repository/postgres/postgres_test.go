package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Developerproject2024/devboard/internal/domain"
	generated "github.com/Developerproject2024/devboard/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

func mustUUID(t *testing.T, raw string) pgtype.UUID {
	t.Helper()
	var id pgtype.UUID
	if err := id.Scan(raw); err != nil {
		t.Fatalf("scan uuid %q: %v", raw, err)
	}
	return id
}

func mustTimestamptz(t *testing.T, ts time.Time) pgtype.Timestamptz {
	t.Helper()
	return pgtype.Timestamptz{Time: ts, Valid: true}
}

func TestToDomainTaskMapsGeneratedTask(t *testing.T) {
	id := mustUUID(t, "11111111-1111-1111-1111-111111111111")
	projectID := mustUUID(t, "22222222-2222-2222-2222-222222222222")
	createdBy := mustUUID(t, "33333333-3333-3333-3333-333333333333")
	assigneeID := mustUUID(t, "44444444-4444-4444-4444-444444444444")
	createdAt := time.Date(2026, time.September, 14, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(5 * time.Minute)

	got := toDomainTask(generated.Task{
		ID:          id,
		ProjectID:   projectID,
		Title:       "Tarea",
		Description: "Descripcion",
		Status:      string(domain.TaskStatusTodo),
		AssigneeID:  assigneeID,
		CreatedBy:   createdBy,
		CreatedAt:   mustTimestamptz(t, createdAt),
		UpdatedAt:   mustTimestamptz(t, updatedAt),
	})

	if got.ID != id.String() {
		t.Fatalf("ID = %q, want %q", got.ID, id.String())
	}
	if got.ProjectID != projectID.String() {
		t.Fatalf("ProjectID = %q, want %q", got.ProjectID, projectID.String())
	}
	if got.Title != "Tarea" || got.Description != "Descripcion" {
		t.Fatalf("title/description not mapped as expected: %+v", got)
	}
	if got.Status != domain.TaskStatusTodo {
		t.Fatalf("status = %q, want %q", got.Status, domain.TaskStatusTodo)
	}
	if got.CreatedBy != createdBy.String() {
		t.Fatalf("CreatedBy = %q, want %q", got.CreatedBy, createdBy.String())
	}
	if got.AssigneeID == nil || *got.AssigneeID != assigneeID.String() {
		t.Fatalf("AssigneeID = %v, want %q", got.AssigneeID, assigneeID.String())
	}
	if !got.CreatedAt.Equal(createdAt) || !got.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("timestamps not mapped: created=%v updated=%v", got.CreatedAt, got.UpdatedAt)
	}
}

func TestToDomainUserMapsGeneratedUser(t *testing.T) {
	id := mustUUID(t, "55555555-5555-5555-5555-555555555555")
	createdAt := time.Date(2026, time.September, 14, 11, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(10 * time.Minute)

	got := toDomainUser(generated.User{
		ID:           id,
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Ana",
		CreatedAt:    mustTimestamptz(t, createdAt),
		UpdatedAt:    mustTimestamptz(t, updatedAt),
	})

	if got.ID != id.String() || got.Email != "test@example.com" || got.PasswordHash != "hash" || got.Name != "Ana" {
		t.Fatalf("user mapping failed: %+v", got)
	}
	if !got.CreatedAt.Equal(createdAt) || !got.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("user timestamps not mapped: created=%v updated=%v", got.CreatedAt, got.UpdatedAt)
	}
}

func TestTaskRepositoryCreateRejectsInvalidIdentifiers(t *testing.T) {
	repo := &TaskRepository{}
	badProjectTask := &domain.Task{ProjectID: "not-a-uuid", CreatedBy: "55555555-5555-5555-5555-555555555555"}
	if err := repo.Create(context.Background(), badProjectTask); err == nil || !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Create should reject invalid ProjectID: %v", err)
	}

	badCreatedByTask := &domain.Task{ProjectID: "55555555-5555-5555-5555-555555555555", CreatedBy: "bad-uuid"}
	if err := repo.Create(context.Background(), badCreatedByTask); err == nil || !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Create should reject invalid CreatedBy: %v", err)
	}
}

func TestTaskRepositoryUpdateRejectsInvalidIdentifiers(t *testing.T) {
	repo := &TaskRepository{}
	badTaskID := &domain.Task{ID: "bad-id", Title: "title", Description: "desc", Status: domain.TaskStatusTodo}
	if err := repo.Update(context.Background(), badTaskID); err == nil || !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Update should reject invalid Task ID: %v", err)
	}

	badAssignee := &domain.Task{ID: "55555555-5555-5555-5555-555555555555", Title: "title", Description: "desc", Status: domain.TaskStatusTodo, AssigneeID: nil}
	badAssignee.AssigneeID = new(string)
	*badAssignee.AssigneeID = "also-bad"
	if err := repo.Update(context.Background(), badAssignee); err == nil || !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Update should reject invalid AssigneeID: %v", err)
	}
}

func TestTaskRepositoryGetByIDRejectsInvalidIdentifier(t *testing.T) {
	repo := &TaskRepository{}
	_, err := repo.GetByID(context.Background(), "not-a-uuid")
	if err == nil || !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("GetByID should reject invalid UUID and return ErrNotFound: %v", err)
	}
}

func TestUserRepositoryGetByIDRejectsInvalidIdentifier(t *testing.T) {
	repo := &UserRepository{}
	_, err := repo.GetByID(context.Background(), "not-a-uuid")
	if err == nil || !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("GetByID should reject invalid UUID and return ErrNotFound: %v", err)
	}
}
