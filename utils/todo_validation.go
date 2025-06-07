package utils

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func IsValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^628\d{9,14}$`)
	return re.MatchString(phone)
}

func IsValidDate(date string) bool {
	re := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
	return re.MatchString(date)
}


func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPassword(password string) bool {

	if len(password) < 8 {
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

func HashingPlainText(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
	}
	return string(bytes), nil
}

func CompareHashAndPlainText(hashedText, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plaintext))
	return err == nil
}

func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

func ValidatePin(pin string) error {
	if len(pin) != 6 {
		return errors.New("PIN must be 6 digits")
	}
	if !IsNumeric(pin) {
		return errors.New("PIN must contain only numbers")
	}
	return nil
}

func GenerateOTP() (string, error) {
	max := big.NewInt(9999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", n.Int64()), nil
}
