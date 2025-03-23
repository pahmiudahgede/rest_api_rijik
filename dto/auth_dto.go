package dto

import (
	"regexp"
	"strings"
)

type RegisterRequest struct {
	RoleID string `json:"role_id"`
	Phone  string `json:"phone"`
}

type VerifyOTPRequest struct {
	RoleID string `json:"role_id"`
	Phone  string `json:"phone"`
	OTP   string `json:"otp"`
}

type UserDataResponse struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	Token    string `json:"token"`
}

func (r *RegisterRequest) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.RoleID) == "" {
		errors["role_id"] = append(errors["role_id"], "Role ID is required")
	}

	if strings.TrimSpace(r.Phone) == "" {
		errors["phone"] = append(errors["phone"], "Phone is required")
	} else if !IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Invalid phone number format. Use 62 followed by 9-13 digits")
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
