package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Developerproject2024/devboard/internal/repository/memory"
	"github.com/Developerproject2024/devboard/internal/usecase"
	"github.com/Developerproject2024/devboard/internal/validator"
)

func newUserHandlerForTest() *UserHandler {
	return NewUserHandler(
		usecase.NewUserUseCase(memory.NewUserRepository()),
		validator.New(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
}

func userRequest(method, path, body string) *http.Request {
	request := httptest.NewRequest(method, path, stringReader(body))
	request.SetPathValue("id", "user_1")
	return request
}

func stringReader(value string) *reader {
	return &reader{value: value}
}

type reader struct {
	value string
}

func (r *reader) Read(buffer []byte) (int, error) {
	if r.value == "" {
		return 0, io.EOF
	}
	n := copy(buffer, r.value)
	r.value = r.value[n:]
	return n, nil
}

func TestUserHandlerCreateAndGet(t *testing.T) {
	handler := newUserHandlerForTest()
	createRecorder := httptest.NewRecorder()

	handler.Create(createRecorder, userRequest(http.MethodPost, "/users", `{"email":"ana@example.com","name":"Ana","password":"password"}`))
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("Create() status = %d; se esperaba %d", createRecorder.Code, http.StatusCreated)
	}

	var created userResponse
	if err := json.NewDecoder(createRecorder.Body).Decode(&created); err != nil {
		t.Fatalf("Create() devolvió JSON inválido: %v", err)
	}

	getRequest := userRequest(http.MethodGet, "/users/"+created.ID, "")
	getRequest.SetPathValue("id", created.ID)
	getRecorder := httptest.NewRecorder()
	handler.Get(getRecorder, getRequest)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("Get() status = %d; se esperaba %d", getRecorder.Code, http.StatusOK)
	}
}

func TestUserHandlerCreateInvalidRequest(t *testing.T) {
	handler := newUserHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Create(recorder, userRequest(http.MethodPost, "/users", `{}`))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Create() status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUserHandlerCreateAlreadyExists(t *testing.T) {
	handler := newUserHandlerForTest()
	body := `{"email":"ana@example.com","name":"Ana","password":"password"}`
	firstRecorder := httptest.NewRecorder()
	handler.Create(firstRecorder, userRequest(http.MethodPost, "/users", body))

	secondRecorder := httptest.NewRecorder()
	handler.Create(secondRecorder, userRequest(http.MethodPost, "/users", body))

	if secondRecorder.Code != http.StatusConflict {
		t.Fatalf("Create() status = %d; se esperaba %d", secondRecorder.Code, http.StatusConflict)
	}
}

func TestUserHandlerCreateMalformedJSON(t *testing.T) {
	handler := newUserHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Create(recorder, userRequest(http.MethodPost, "/users", `{`))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Create() status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUserHandlerGetNotFound(t *testing.T) {
	handler := newUserHandlerForTest()
	recorder := httptest.NewRecorder()

	handler.Get(recorder, userRequest(http.MethodGet, "/users/user_1", ""))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("Get() status = %d; se esperaba %d", recorder.Code, http.StatusNotFound)
	}
}
