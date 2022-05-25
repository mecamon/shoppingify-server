//go:build !integration
// +build !integration

package utils

import "testing"

func TestGenerateHash(t *testing.T) {
	pass := "PasswordValid890"
	hashedPassword, err := GenerateHash(pass)

	if err != nil || len(hashedPassword) <= len(pass) {
		t.Errorf("could not hash the password: %v", pass)
	}
}

var compareHashAndPassTests = []struct {
	password      string
	expectedValue bool
}{
	{password: "PasswordValid123456", expectedValue: true},
	{password: "", expectedValue: false},
}

func TestCompareHashAndPass(t *testing.T) {
	for _, tt := range compareHashAndPassTests {
		hashedPassword, _ := GenerateHash(tt.password)
		isCorrect, _ := CompareHashAndPass(hashedPassword, tt.password)

		if !isCorrect {
			t.Errorf("got %v when expected value was %v", isCorrect, tt.expectedValue)
		}
	}
}
