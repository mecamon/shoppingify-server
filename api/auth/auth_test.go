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
		auth:           models.User{Email: "invalid@!mail", Password: "invalid-pass", Username: "cac"},
		expectedValid:  false,
		expectedErrors: 3,
	},
	{
		name:           "valid credentials",
		auth:           models.User{Email: "valid@email.com", Password: "SecurePassword1123", Username: "Bonjovi"},
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
