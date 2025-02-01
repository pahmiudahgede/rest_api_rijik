package dto

import (
	"regexp"
	"strings"
)

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

type UpdateUserDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (r *UpdateUserDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

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

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func IsUpdateValidPhoneNumber(phone string) bool {

	re := regexp.MustCompile(`^\+62\d{9,13}$`)
	return re.MatchString(phone)
}

func IsUPdateValidEmail(email string) bool {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

type UpdatePasswordDTO struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

func (u *UpdatePasswordDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if u.OldPassword == "" {
		errors["old_password"] = append(errors["old_password"], "Old password is required")
	}

	if u.NewPassword == "" {
		errors["new_password"] = append(errors["new_password"], "New password is required")
	} else if len(u.NewPassword) < 8 {
		errors["new_password"] = append(errors["new_password"], "Password must be at least 8 characters long")
	}

	if u.ConfirmNewPassword == "" {
		errors["confirm_new_password"] = append(errors["confirm_new_password"], "Confirm new password is required")
	} else if u.NewPassword != u.ConfirmNewPassword {
		errors["confirm_new_password"] = append(errors["confirm_new_password"], "Passwords do not match")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
