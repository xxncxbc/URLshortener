package auth

import (
	"URLshortener/internal/user"
	"testing"
)

type MockUserRepository struct {
}

func (mockRepo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "test@example.com",
	}, nil
}

func (mockRepo *MockUserRepository) GetByEmail(e string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const InitialEmail = "test@example.com"
	authService := NewAuthService(&MockUserRepository{})
	email, _, err := authService.Register(InitialEmail, "", "testservice")
	if err != nil {
		t.Fatal(err)
		return
	}
	if email != InitialEmail {
		t.Fatalf("Email does not match. Expected %s, got %s", InitialEmail, email)
		return
	}

}
