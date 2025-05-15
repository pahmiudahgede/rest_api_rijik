package utils

import (
	"fmt"
	"regexp"
	"strconv"
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

func ValidateFloatPrice(price string) (float64, error) {

	// price = strings.Trim(price, `"`)
	// price = strings.TrimSpace(price)

	parsedPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, fmt.Errorf("harga tidak valid. Format harga harus angka desimal.")
	}

	if parsedPrice <= 0 {
		return 0, fmt.Errorf("harga harus lebih besar dari 0.")
	}

	return parsedPrice, nil
}
