package userpin

import (
	"rijig/utils"
	"strings"
)

type RequestPinDTO struct {
	// DeviceId string `json:"device_id"`
	Pin      string `json:"userpin"`
}

func (r *RequestPinDTO) ValidateRequestPinDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if err := utils.ValidatePin(r.Pin); err != nil {
		errors["pin"] = append(errors["pin"], err.Error())
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

type UpdatePinDTO struct {
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

func (u *UpdatePinDTO) ValidateUpdatePinDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(u.OldPin) == "" {
		errors["old_pin"] = append(errors["old_pin"], "Old pin is required")
	}

	if err := utils.ValidatePin(u.NewPin); err != nil {
		errors["new_pin"] = append(errors["new_pin"], err.Error())
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
