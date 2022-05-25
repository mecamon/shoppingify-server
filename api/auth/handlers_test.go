//go:build integration
// +build integration

package auth

import (
	"bytes"
	"encoding/json"
	"github.com/mecamon/shoppingify-server/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

var registerTests = []struct {
	testName           string
	user               models.User
	expectedStatusCode int
}{
	{testName: "invalid-user-data", user: models.User{
		Name:     "c",
		Lastname: "a",
		Email:    "mail-not-valid.com",
		Password: "1234566",
	}, expectedStatusCode: http.StatusBadRequest},
	{testName: "email-in-use", user: models.User{
		Name:     "Pepe",
		Lastname: "Pepega",
		Email:    LoginUserData.Email,
		Password: LoginUserData.Password,
	}, expectedStatusCode: http.StatusConflict},
	{testName: "valid-user-data", user: models.User{
		Name:     "Carlos",
		Lastname: "Mejia",
		Email:    "carlos@mejia.com",
		Password: "ValidPassword1233",
	}, expectedStatusCode: http.StatusCreated},
}

func TestHandler_Register(t *testing.T) {
	for _, tt := range registerTests {
		t.Log(tt.testName)
		{
			body, _ := json.Marshal(&tt.user)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
			Router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Got %v statusCode when expected %v", rr.Code, tt.expectedStatusCode)
			}
		}
	}
}

var loginTests = []struct {
	testName           string
	email              string
	password           string
	expectedStatusCode int
}{
	{testName: "valid-credentials", email: LoginUserData.Email, password: LoginUserData.Password, expectedStatusCode: http.StatusOK},
	{testName: "invalid-email", email: "user@notloged.com", password: LoginUserData.Password, expectedStatusCode: http.StatusBadRequest},
	{testName: "invalid-email", email: "user@notloged.com", password: LoginUserData.Password, expectedStatusCode: http.StatusBadRequest},
	{testName: "invalid-email", email: "user@notloged.com", password: LoginUserData.Password, expectedStatusCode: http.StatusBadRequest},
	{testName: "invalid-password", email: "charles@mail.com", password: "wrongPass", expectedStatusCode: http.StatusBadRequest},
}

func TestHandler_Login(t *testing.T) {
	for _, tt := range loginTests {
		t.Log(tt.testName)
		auth := models.Auth{Email: tt.email, Password: tt.password}
		body, _ := json.Marshal(&auth)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		Router.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("got %v when expected was %v", rr.Code, tt.expectedStatusCode)
			t.Error("Error:", rr.Body)
		}
	}
}
