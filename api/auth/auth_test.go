//go:build !integration
// +build !integration

package auth

import (
	"github.com/mecamon/shoppingify-server/models"
	"testing"
)

var validCredTests = []struct {
	name           string
	auth           models.User
	expectedValid  bool
	expectedErrors int
}{
	{
		name:           "invalid credentials 1",
		auth:           models.User{Email: "invalid@!mail", Password: "invalid-pass", Name: "c", Lastname: "d"},
		expectedValid:  false,
		expectedErrors: 4,
	},
	{
		name:           "valid credentials",
		auth:           models.User{Email: "valid@email.com", Password: "SecurePassword1123", Name: "Bonjovi", Lastname: "Jon"},
		expectedValid:  true,
		expectedErrors: 0,
	},
}

func TestValidCredentials(t *testing.T) {
	for _, tt := range validCredTests {
		t.Log(tt.name)
		valid, errMap := validCredentials(tt.auth, "en-EN")
		if valid != tt.expectedValid {
			t.Errorf("expected %v but got %v", tt.expectedValid, valid)
		}
		if !valid && len(errMap) != tt.expectedErrors {
			t.Errorf("Expected %v errors but got %v", tt.expectedErrors, len(errMap))
		}
	}
}

func TestCompleteUserInformation(t *testing.T) {
	user := models.User{
		Name:     "Charles",
		Lastname: "Maze",
		Email:    "charles@maze.com",
		Password: "AValidPass34",
	}
	completedUser, _ := completeUserInformation(user)

	if completedUser.Password == user.Password {
		t.Error("Password is not hashed")
	}
	if completedUser.IsActive == false {
		t.Error("IsActive prop is false")
	}
	if completedUser.LoginCode == "" {
		t.Error("LoginCode prop is empty")
	}
	if completedUser.CreatedAt == 0 || completedUser.UpdatedAt == 0 {
		t.Error("CreatedAt and/or UpdatedAt unset")
	}
}

func TestCreateVisitorInformation(t *testing.T) {
	visitor := createVisitorInformation()

	if !visitor.IsVisitor {
		t.Error("IsVisitor is false")
	}
	if visitor.LoginCode == "" {
		t.Error("LoginCode prop is empty")
	}
	if visitor.CreatedAt == 0 || visitor.UpdatedAt == 0 {
		t.Error("CreatedAt and/or UpdatedAt unset")
	}
}
