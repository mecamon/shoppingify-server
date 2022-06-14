//go:build !integration
// +build !integration

package json_web_token

import "testing"

func TestGenerate(t *testing.T) {
	var id int64 = 12345
	email := "some@mail.com"

	token, err := Generate(id, email)
	if err != nil {
		t.Error(err.Error())
	}
	if token == "" {
		t.Error("error generating token")
	}
}

func TestValidate(t *testing.T) {
	var id int64 = 12345
	email := "some@mail.com"
	token, _ := Generate(id, email)
	claims, err := Validate(token)

	if err != nil {
		t.Error(err.Error())
	}
	if claims.ID != id {
		t.Errorf("expected id was %d but got %d", id, claims.ID)
	}
}
