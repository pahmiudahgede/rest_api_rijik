package dto

import (
	"regexp"
	"strings"
)

type LoginDTO struct {
	RoleID     string `json:"roleid"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type UserResponseWithToken struct {
	UserID   string `json:"user_id"`
	RoleName string `json:"loginas"`
	Token    string `json:"token"`
}

type RegisterDTO struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	RoleID          string `json:"roleId,omitempty"`
}

type UserResponseDTO struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	RoleName      string `json:"role"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func (l *LoginDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(l.RoleID) == "" {
		errors["roleid"] = append(errors["roleid"], "Role ID is required")
	}
	if strings.TrimSpace(l.Identifier) == "" {
		errors["identifier"] = append(errors["identifier"], "Identifier (username, email, or phone) is required")
	}
	if strings.TrimSpace(l.Password) == "" {
		errors["password"] = append(errors["password"], "Password is required")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

func (r *RegisterDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Username) == "" {
		errors["username"] = append(errors["username"], "Username is required")
	}
	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	}

	if strings.TrimSpace(r.Phone) == "" {
		errors["phone"] = append(errors["phone"], "Phone number is required")
	} else if !IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Invalid phone number format. Use +62 followed by 9-13 digits")
	}

	if strings.TrimSpace(r.Email) == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else if !IsValidEmail(r.Email) {
		errors["email"] = append(errors["email"], "Invalid email format")
	}

	if strings.TrimSpace(r.Password) == "" {
		errors["password"] = append(errors["password"], "Password is required")
	} else if !IsValidPassword(r.Password) {
		errors["password"] = append(errors["password"], "Password must be at least 8 characters long and contain at least one number")
	}

	if strings.TrimSpace(r.ConfirmPassword) == "" {
		errors["confirm_password"] = append(errors["confirm_password"], "Confirm password is required")
	} else if r.Password != r.ConfirmPassword {
		errors["confirm_password"] = append(errors["confirm_password"], "Password and confirm password do not match")
	}
	if strings.TrimSpace(r.RoleID) == "" {
		errors["roleId"] = append(errors["roleId"], "RoleID is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func IsValidPhoneNumber(phone string) bool {

	re := regexp.MustCompile(`^\+62\d{9,13}$`)
	return re.MatchString(phone)
}

func IsValidEmail(email string) bool {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	re := regexp.MustCompile(`\d`)
	return re.MatchString(password)
}
