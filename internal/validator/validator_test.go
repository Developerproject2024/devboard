package validator

import (
	"testing"
)

type TestStruct struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=3,max=50"`
	Age      int    `validate:"required,min=18"`
	Username string `validate:"required,min=5,max=20"`
}

func TestNewValidator(t *testing.T) {
	v := New()
	if v == nil {
		t.Fatalf("esperaba un validator válido, recibí nil")
	}
	if v.validate == nil {
		t.Fatalf("esperaba validator.Validate inicializado, recibí nil")
	}
}

func TestValidateWithValidStruct(t *testing.T) {
	v := New()
	s := TestStruct{
		Email:    "test@example.com",
		Name:     "John",
		Age:      25,
		Username: "john_doe",
	}

	errors := v.Validate(s)
	if len(errors) > 0 {
		t.Fatalf("esperaba 0 errores, recibí %d: %v", len(errors), errors)
	}
}

func TestValidateWithMissingRequiredFields(t *testing.T) {
	v := New()
	s := TestStruct{
		Email:    "",
		Name:     "",
		Age:      0,
		Username: "",
	}

	errors := v.Validate(s)
	if len(errors) == 0 {
		t.Fatalf("esperaba errores, pero recibí ninguno")
	}

	// Should have errors for each required field
	if len(errors) < 3 {
		t.Fatalf("esperaba al menos 3 errores, recibí %d", len(errors))
	}

	// Verify that required field errors are present
	requiredFieldsFound := 0
	for _, err := range errors {
		if err.Message == "esta campo es obligatorio" {
			requiredFieldsFound++
		}
	}

	if requiredFieldsFound == 0 {
		t.Fatalf("esperaba errores de campo obligatorio, pero no los encontré")
	}
}

func TestValidateWithInvalidEmail(t *testing.T) {
	v := New()
	s := TestStruct{
		Email:    "not-an-email",
		Name:     "John",
		Age:      25,
		Username: "john_doe",
	}

	errors := v.Validate(s)
	if len(errors) == 0 {
		t.Fatalf("esperaba errores de email, pero recibí ninguno")
	}

	// Find email error
	emailErrorFound := false
	for _, err := range errors {
		if err.Field == "email" && err.Message == "debe ser un correo válido" {
			emailErrorFound = true
			break
		}
	}

	if !emailErrorFound {
		t.Fatalf("esperaba error de email inválido")
	}
}

func TestValidateWithMinLengthViolation(t *testing.T) {
	v := New()
	s := TestStruct{
		Email:    "test@example.com",
		Name:     "Jo", // Less than min=3
		Age:      25,
		Username: "john_doe",
	}

	errors := v.Validate(s)
	if len(errors) == 0 {
		t.Fatalf("esperaba errores de min length, pero recibí ninguno")
	}

	// Find name error
	nameErrorFound := false
	for _, err := range errors {
		if err.Field == "name" && err.Message == "debe tener mínimo 3 caracteres" {
			nameErrorFound = true
			break
		}
	}

	if !nameErrorFound {
		t.Fatalf("esperaba error de min length en nombre")
	}
}

func TestValidateWithMaxLengthViolation(t *testing.T) {
	v := New()
	s := TestStruct{
		Email:    "test@example.com",
		Name:     "John",
		Age:      25,
		Username: "this_is_a_very_long_username_that_exceeds_max_length",
	}

	errors := v.Validate(s)
	if len(errors) == 0 {
		t.Fatalf("esperaba errores de max length, pero recibí ninguno")
	}

	// Find username error
	usernameErrorFound := false
	for _, err := range errors {
		if err.Field == "username" && err.Message == "debe tener máximo 20 caracteres" {
			usernameErrorFound = true
			break
		}
	}

	if !usernameErrorFound {
		t.Fatalf("esperaba error de max length en username")
	}
}

func TestValidateWithMultipleErrors(t *testing.T) {
	v := New()
	s := TestStruct{
		Email:    "invalid-email",
		Name:     "Jo",
		Age:      15, // Less than min=18
		Username: "ab",
	}

	errors := v.Validate(s)
	if len(errors) == 0 {
		t.Fatalf("esperaba múltiples errores, pero recibí ninguno")
	}

	if len(errors) < 3 {
		t.Fatalf("esperaba al menos 3 errores, recibí %d", len(errors))
	}
}

func TestValidateWithNilInput(t *testing.T) {
	v := New()
	errors := v.Validate(nil)
	if errors != nil {
		t.Fatalf("esperaba nil para entrada nil, recibí: %v", errors)
	}
}

func TestTranslateTagRequired(t *testing.T) {
	// This is tested indirectly through Validate, but we can test the behavior
	v := New()
	s := TestStruct{
		Email: "",
		Name:  "",
		Age:   0,
	}

	errors := v.Validate(s)
	if errors == nil {
		t.Fatalf("esperaba errores")
	}

	for _, err := range errors {
		if err.Message == "" {
			t.Fatalf("esperaba mensaje de error no vacío")
		}
	}
}

func TestValidateFieldNameLowercase(t *testing.T) {
	v := New()
	s := TestStruct{
		Email: "invalid",
	}

	errors := v.Validate(s)
	if len(errors) == 0 {
		t.Fatalf("esperaba errores")
	}

	// Check that field names are lowercase
	for _, err := range errors {
		if err.Field != "" && err.Field[0] >= 'A' && err.Field[0] <= 'Z' {
			t.Fatalf("esperaba nombre de campo en minúsculas, recibí: %s", err.Field)
		}
	}
}

func TestValidateWithEdgeCaseValues(t *testing.T) {
	v := New()

	// Test with minimum valid values
	s := TestStruct{
		Email:    "a@b.c",
		Name:     "Bob",
		Age:      18,
		Username: "alice",
	}

	errors := v.Validate(s)
	if len(errors) > 0 {
		t.Fatalf("esperaba 0 errores para valores válidos, recibí: %v", errors)
	}

	// Test with maximum valid values
	s2 := TestStruct{
		Email:    "verylongemailaddress@verylongdomainname.com",
		Name:     "ABC",
		Age:      150,
		Username: "abcde",
	}

	errors2 := v.Validate(s2)
	if len(errors2) > 0 {
		t.Fatalf("esperaba 0 errores para valores válidos, recibí: %v", errors2)
	}
}
