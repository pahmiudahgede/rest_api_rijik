package utils

import (
	"regexp"
	"strings"
)

func IsValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^62\d{9,13}$`)
	return re.MatchString(phone)
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPassword(password string) bool {

	if len(password) < 6 {
		return false
	}

	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else if isSpecialCharacter(char) {
			hasSpecial = true
		}
	}

	return hasUpper && hasDigit && hasSpecial
}

func isSpecialCharacter(char rune) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"
	return strings.ContainsRune(specialChars, char)
}
