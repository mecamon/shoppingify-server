//go:build !integration
// +build !integration

package utils

import "testing"

var validEmailTests = []struct {
	name     string
	email    string
	expected bool
}{
	{name: "invalid email 1", email: "not-valid-email", expected: false},
	{name: "invalid email 2", email: "not-valid@email", expected: false},
	{name: "invalid email 3", email: "not-valid.email", expected: false},
	{name: "valid email", email: "valid@mail.com", expected: true},
}

func TestHasValidEmail(t *testing.T) {
	for _, tt := range validEmailTests {
		t.Logf(tt.name)
		{
			matched := HasValidEmail(tt.email)
			if matched != tt.expected {
				t.Errorf("Got %v when expected was %v", matched, tt.expected)
			}
		}
	}
}

var validPassTests = []struct {
	name     string
	password string
	expected bool
}{
	{name: "invalid password 1", password: "password", expected: false},
	{name: "invalid password 2", password: "pAssworD", expected: false},
	{name: "invalid password 3", password: "PASSWORD1", expected: false},
	{name: "valid password 1", password: "PassWord12", expected: true},
	{name: "valid password 2", password: "PASSWORd12", expected: true},
}

func TestValidPassword(t *testing.T) {
	for _, tt := range validPassTests {
		t.Log(tt.name)
		{
			if matched := HasValidPass(tt.password); matched != tt.expected {
				t.Errorf("Got %v when expected %v", matched, tt.expected)
			}
		}
	}
}
