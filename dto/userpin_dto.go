package dto

import (
	"fmt"
	"regexp"
	"strings"
)

type RequestUserPinDTO struct {
	Pin string `json:"userpin"`
}

func (r *RequestUserPinDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Pin) == "" {
		errors["pin"] = append(errors["pin"], "Pin is required")
	}

	if err := validatePin(r.Pin); err != nil {
		errors["pin"] = append(errors["pin"], err.Error())
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

func (u *UpdateUserPinDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(u.OldPin) == "" {
		errors["old_pin"] = append(errors["old_pin"], "Old pin is required")
	}

	if err := validatePin(u.NewPin); err != nil {
		errors["new_pin"] = append(errors["new_pin"], err.Error())
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

func validatePin(pin string) error {
	if len(pin) != 6 {
		return fmt.Errorf("pin harus terdiri dari 6 digit")
	} else if !isNumeric(pin) {
		return fmt.Errorf("pin harus berupa angka")
	}
	return nil
}
