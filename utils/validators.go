package utils

import "regexp"

func HasValidEmail(email string) bool {
	pattern := "^.{2,}@.{2,}\\..{2,}$"
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return matched
}

func HasValidPass(password string) bool {
	patterns := []string{".*[a-z]", ".*[A-Z]", ".*\\d", "^.{8,}$"}
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, password)
		if !matched || err != nil {
			return false
		}
	}
	return true
}
