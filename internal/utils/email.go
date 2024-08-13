package utils

import (
	"net/mail"
	"regexp"
)

func IsValidEmail(email string) bool {
	var err error

	// Regular expression pattern for validating email addresses
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	_, err = mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Check if the email matches the pattern
	return regex.MatchString(email) && err == nil
}
