package dto

import (
	"regexp"
	"strings"
)

type LoginAdminRequest struct {
	Deviceid string `json:"device_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

type RegisterAdminRequest struct {
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	Dateofbirth     string `json:"dateofbirth"`
	Placeofbirth    string `json:"placeofbirth"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserAdminDataResponse struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	Dateofbirth  string `json:"dateofbirth"`
	Placeofbirth string `json:"placeofbirth"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func (r *RegisterAdminRequest) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	}

	if strings.TrimSpace(r.Gender) == "" {
		errors["gender"] = append(errors["gender"], "Gender is required")
	} else if r.Gender != "male" && r.Gender != "female" {
		errors["gender"] = append(errors["gender"], "Gender must be either 'male' or 'female'")
	}

	if strings.TrimSpace(r.Dateofbirth) == "" {
		errors["dateofbirth"] = append(errors["dateofbirth"], "Date of birth is required")
	}

	if strings.TrimSpace(r.Placeofbirth) == "" {
		errors["placeofbirth"] = append(errors["placeofbirth"], "Place of birth is required")
	}

	if strings.TrimSpace(r.Phone) == "" {
		errors["phone"] = append(errors["phone"], "Phone is required")
	} else if !IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Invalid phone number format. Use 62 followed by 9-13 digits")
	}

	if strings.TrimSpace(r.Email) == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else if !IsValidEmail(r.Email) {
		errors["email"] = append(errors["email"], "Invalid email format")
	}

	if len(r.Password) < 6 {
		errors["password"] = append(errors["password"], "Password must be at least 6 characters")
	} else if !IsValidPassword(r.Password) {
		errors["password"] = append(errors["password"], "Password must contain at least one uppercase letter, one number, and one special character")
	}

	if r.Password != r.PasswordConfirm {
		errors["password_confirm"] = append(errors["password_confirm"], "Password and confirmation do not match")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

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
