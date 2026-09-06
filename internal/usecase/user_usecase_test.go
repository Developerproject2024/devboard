package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Developerproject2024/devboard/internal/domain"
)

type userRepositoryStub struct {
	user        *domain.User
	byEmailErr  error
	byIDErr     error
	createErr   error
	createdUser *domain.User
}

func (r *userRepositoryStub) Create(_ context.Context, user *domain.User) error {
	r.createdUser = user
	return r.createErr
}

func (r *userRepositoryStub) GetByID(_ context.Context, _ string) (*domain.User, error) {
	return r.user, r.byIDErr
}

func (r *userRepositoryStub) GetByEmail(_ context.Context, _ string) (*domain.User, error) {
	return r.user, r.byEmailErr
}

func TestUserUseCase_CreateUser(t *testing.T) {
	repository := &userRepositoryStub{byEmailErr: domain.ErrNotFound}
	useCase := NewUserUseCase(repository)

	user, err := useCase.CreateUser(context.Background(), "ana@example.com", "Ana", "hash")
	if err != nil {
		t.Fatalf("CreateUser() devolvió un error inesperado: %v", err)
	}
	if user.Email != "ana@example.com" || repository.createdUser != user {
		t.Fatalf("CreateUser() devolvió un usuario inesperado: %+v", user)
	}
}

func TestUserUseCase_CreateUserAlreadyExists(t *testing.T) {
	existing := domain.NewUser("ana@example.com", "Ana", "hash")
	useCase := NewUserUseCase(&userRepositoryStub{user: existing})

	_, err := useCase.CreateUser(context.Background(), existing.Email, existing.Name, existing.PasswordHash)
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("CreateUser() devolvió %v; se esperaba ErrAlreadyExists", err)
	}
}

func TestUserUseCase_CreateUserCheckEmailError(t *testing.T) {
	checkErr := errors.New("fallo consultando email")
	useCase := NewUserUseCase(&userRepositoryStub{byEmailErr: checkErr})

	_, err := useCase.CreateUser(context.Background(), "ana@example.com", "Ana", "hash")
	if !errors.Is(err, checkErr) {
		t.Fatalf("CreateUser() devolvió %v; se esperaba %v", err, checkErr)
	}
}

func TestUserUseCase_CreateUserInvalid(t *testing.T) {
	useCase := NewUserUseCase(&userRepositoryStub{byEmailErr: domain.ErrNotFound})

	_, err := useCase.CreateUser(context.Background(), "", "Ana", "hash")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("CreateUser() devolvió %v; se esperaba ErrInvalidInput", err)
	}
}

func TestUserUseCase_CreateUserRepositoryError(t *testing.T) {
	createErr := errors.New("fallo de persistencia")
	useCase := NewUserUseCase(&userRepositoryStub{byEmailErr: domain.ErrNotFound, createErr: createErr})

	_, err := useCase.CreateUser(context.Background(), "ana@example.com", "Ana", "hash")
	if !errors.Is(err, createErr) {
		t.Fatalf("CreateUser() devolvió %v; se esperaba %v", err, createErr)
	}
}

func TestUserUseCase_GetUser(t *testing.T) {
	user := domain.NewUser("ana@example.com", "Ana", "hash")
	useCase := NewUserUseCase(&userRepositoryStub{user: user})

	found, err := useCase.GetUser(context.Background(), "user-1")
	if err != nil || found != user {
		t.Fatalf("GetUser() devolvió usuario=%+v, error=%v", found, err)
	}
}

func TestUserUseCase_GetUserError(t *testing.T) {
	getErr := errors.New("usuario no encontrado")
	useCase := NewUserUseCase(&userRepositoryStub{byIDErr: getErr})

	_, err := useCase.GetUser(context.Background(), "user-1")
	if !errors.Is(err, getErr) {
		t.Fatalf("GetUser() devolvió %v; se esperaba %v", err, getErr)
	}
}
