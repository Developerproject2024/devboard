package domain

import "testing"

func TestNewUser(t *testing.T) {
	user := NewUser("ana@example.com", "Ana", "hash")

	if user.Email != "ana@example.com" || user.Name != "Ana" || user.PasswordHash != "hash" {
		t.Fatalf("NewUser() devolvió un usuario inesperado: %+v", user)
	}
	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Fatalf("NewUser() no estableció timestamps: %+v", user)
	}
}

func TestUserIsValid(t *testing.T) {
	tests := []struct {
		name  string
		user  *User
		valid bool
	}{
		{name: "valid", user: &User{Email: "ana@example.com", Name: "Ana"}, valid: true},
		{name: "missing email", user: &User{Name: "Ana"}},
		{name: "missing name", user: &User{Email: "ana@example.com"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.user.IsValid(); got != test.valid {
				t.Fatalf("IsValid() = %v; se esperaba %v", got, test.valid)
			}
		})
	}
}
