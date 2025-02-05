package dto

import (
	"regexp"
	"strings"
)

type UserPinResponseDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Pin       string `json:"userpin"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type RequestUserPinDTO struct {
	Pin string `json:"userpin"`
}

func (r *RequestUserPinDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Pin) == "" {
		errors["pin"] = append(errors["pin"], "Pin is required")
	}

	if len(r.Pin) != 6 {
		errors["pin"] = append(errors["pin"], "Pin harus terdiri dari 6 digit")
	} else if !isNumeric(r.Pin) {
		errors["pin"] = append(errors["pin"], "Pin harus berupa angka")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

type UpdateUserPinDTO struct {
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

func (r *UpdateUserPinDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.OldPin) == "" {
		errors["old_pin"] = append(errors["old_pin"], "Old pin is required")
	}

	if strings.TrimSpace(r.NewPin) == "" {
		errors["new_pin"] = append(errors["new_pin"], "New pin is required")
	} else if len(r.NewPin) < 6 {
		errors["new_pin"] = append(errors["new_pin"], "New pin must be at least 6 digits")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func isNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}
