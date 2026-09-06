package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Developerproject2024/devboard/internal/domain"
)

type taskRepositoryStub struct {
	task        *domain.Task
	tasks       []*domain.Task
	createErr   error
	getErr      error
	listErr     error
	updateErr   error
	createdTask *domain.Task
	updatedTask *domain.Task
}

func (r *taskRepositoryStub) Create(_ context.Context, task *domain.Task) error {
	r.createdTask = task
	return r.createErr
}

func (r *taskRepositoryStub) Update(_ context.Context, task *domain.Task) error {
	r.updatedTask = task
	return r.updateErr
}

func (r *taskRepositoryStub) GetByID(_ context.Context, _ string) (*domain.Task, error) {
	return r.task, r.getErr
}

func (r *taskRepositoryStub) ListByProject(_ context.Context, _ string) ([]*domain.Task, error) {
	return r.tasks, r.listErr
}

func newTaskUseCaseStub(repository *taskRepositoryStub) *TaskUseCase {
	return NewTaskUseCase(repository, repository)
}

func TestTaskUseCase_CreateTask(t *testing.T) {
	repository := &taskRepositoryStub{}
	useCase := newTaskUseCaseStub(repository)

	task, err := useCase.CreateTask(context.Background(), "project-1", "Título", "Descripción", "user-1")
	if err != nil {
		t.Fatalf("CreateTask() devolvió un error inesperado: %v", err)
	}
	if task.ProjectID != "project-1" || task.Title != "Título" || task.CreatedBy != "user-1" {
		t.Fatalf("CreateTask() creó una tarea inesperada: %+v", task)
	}
	if repository.createdTask != task {
		t.Fatal("CreateTask() no envió la tarea al repositorio")
	}
}

func TestTaskUseCase_CreateTask_Error(t *testing.T) {
	createErr := errors.New("fallo de persistencia")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{createErr: createErr})

	_, err := useCase.CreateTask(context.Background(), "project-1", "Título", "", "user-1")
	if !errors.Is(err, createErr) {
		t.Fatalf("CreateTask() devolvió %v; se esperaba %v", err, createErr)
	}
}

func TestTaskUseCase_ListProjectTasks(t *testing.T) {
	tasks := []*domain.Task{domain.NewTask("project-1", "Tarea", "", "user-1")}
	useCase := newTaskUseCaseStub(&taskRepositoryStub{tasks: tasks})

	found, err := useCase.ListProjectTasks(context.Background(), "project-1")
	if err != nil {
		t.Fatalf("ListProjectTasks() devolvió un error inesperado: %v", err)
	}
	if len(found) != 1 || found[0] != tasks[0] {
		t.Fatalf("ListProjectTasks() devolvió un resultado inesperado: %+v", found)
	}
}

func TestTaskUseCase_ListProjectTasks_Error(t *testing.T) {
	listErr := errors.New("fallo de consulta")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{listErr: listErr})

	_, err := useCase.ListProjectTasks(context.Background(), "project-1")
	if !errors.Is(err, listErr) {
		t.Fatalf("ListProjectTasks() devolvió %v; se esperaba %v", err, listErr)
	}
}

func TestTaskUseCase_GetTask(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{task: task})

	found, err := useCase.GetTask(context.Background(), "task-1")
	if err != nil {
		t.Fatalf("GetTask() devolvió un error inesperado: %v", err)
	}
	if found != task {
		t.Fatal("GetTask() no devolvió la tarea esperada")
	}
}

func TestTaskUseCase_GetTask_Error(t *testing.T) {
	getErr := errors.New("tarea no encontrada")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{getErr: getErr})

	_, err := useCase.GetTask(context.Background(), "task-1")
	if !errors.Is(err, getErr) {
		t.Fatalf("GetTask() devolvió %v; se esperaba %v", err, getErr)
	}
}

func TestTaskUseCase_AssignTask(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	repository := &taskRepositoryStub{task: task}
	useCase := newTaskUseCaseStub(repository)

	assigned, err := useCase.AssignTask(context.Background(), "task-1", "user-2")
	if err != nil {
		t.Fatalf("AssignTask() devolvió un error inesperado: %v", err)
	}
	if assigned.AssigneeID == nil || *assigned.AssigneeID != "user-2" {
		t.Fatalf("AssignTask() no asignó el usuario esperado: %+v", assigned.AssigneeID)
	}
	if repository.updatedTask != task {
		t.Fatal("AssignTask() no actualizó la tarea en el repositorio")
	}
}

func TestTaskUseCase_AssignTask_UpdateError(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	updateErr := errors.New("fallo de actualización")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{task: task, updateErr: updateErr})

	_, err := useCase.AssignTask(context.Background(), "task-1", "user-2")
	if !errors.Is(err, updateErr) {
		t.Fatalf("AssignTask() devolvió %v; se esperaba %v", err, updateErr)
	}
}

func TestTaskUseCase_AssignTask_GetError(t *testing.T) {
	getErr := errors.New("fallo de consulta")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{getErr: getErr})

	_, err := useCase.AssignTask(context.Background(), "task-1", "user-2")
	if !errors.Is(err, getErr) {
		t.Fatalf("AssignTask() devolvió %v; se esperaba %v", err, getErr)
	}
}

func TestTaskUseCase_AssignTaskDone(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	task.Status = domain.TaskStatusDone
	useCase := newTaskUseCaseStub(&taskRepositoryStub{task: task})

	_, err := useCase.AssignTask(context.Background(), "task-1", "user-2")
	if !errors.Is(err, domain.ErrTaskNotAssignable) {
		t.Fatalf("AssignTask() devolvió %v; se esperaba ErrTaskNotAssignable", err)
	}
}

func TestTaskUseCase_UpdateTaskStatus(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	repository := &taskRepositoryStub{task: task}
	useCase := newTaskUseCaseStub(repository)

	updated, err := useCase.UpdateTaskStatus(context.Background(), "task-1", domain.TaskStatusInProgress)
	if err != nil {
		t.Fatalf("UpdateTaskStatus() devolvió un error inesperado: %v", err)
	}
	if updated.Status != domain.TaskStatusInProgress {
		t.Fatalf("UpdateTaskStatus() dejó el estado %q", updated.Status)
	}
	if repository.updatedTask != task {
		t.Fatal("UpdateTaskStatus() no actualizó la tarea en el repositorio")
	}
}

func TestTaskUseCase_UpdateTaskStatus_InvalidStatus(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	repository := &taskRepositoryStub{task: task}
	useCase := newTaskUseCaseStub(repository)

	_, err := useCase.UpdateTaskStatus(context.Background(), "task-1", domain.TaskStatus("invalid"))
	if err == nil {
		t.Fatal("UpdateTaskStatus() debía rechazar un estado inválido")
	}
	if repository.updatedTask != nil {
		t.Fatal("UpdateTaskStatus() no debía actualizar una tarea con estado inválido")
	}
}

func TestTaskUseCase_UpdateTaskStatus_GetError(t *testing.T) {
	getErr := errors.New("fallo de consulta")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{getErr: getErr})

	_, err := useCase.UpdateTaskStatus(context.Background(), "task-1", domain.TaskStatusInProgress)
	if !errors.Is(err, getErr) {
		t.Fatalf("UpdateTaskStatus() devolvió %v; se esperaba %v", err, getErr)
	}
}

func TestTaskUseCase_UpdateTaskStatus_UpdateError(t *testing.T) {
	task := domain.NewTask("project-1", "Tarea", "", "user-1")
	updateErr := errors.New("fallo de actualización")
	useCase := newTaskUseCaseStub(&taskRepositoryStub{task: task, updateErr: updateErr})

	_, err := useCase.UpdateTaskStatus(context.Background(), "task-1", domain.TaskStatusInProgress)
	if !errors.Is(err, updateErr) {
		t.Fatalf("UpdateTaskStatus() devolvió %v; se esperaba %v", err, updateErr)
	}
}
