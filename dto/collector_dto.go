package dto

import "strings"

type RequestCollectorDTO struct {
	UserId    string `json:"user_id"`
	AddressId string `json:"address_id"`
}

type SelectCollectorRequest struct {
	Collector_id string `json:"collector_id"`
}

func (r *SelectCollectorRequest) ValidateSelectCollectorRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Collector_id) == "" {
		errors["collector_id"] = append(errors["collector_id"], "collector_id harus diisi")
	}
	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

type ResponseCollectorDTO struct {
	ID        string               `json:"collector_id"`
	UserId    string               `json:"user_id"`
	User      []UserResponseDTO    `json:"user,omitempty"`
	AddressId string               `json:"address_id"`
	Address   []AddressResponseDTO `json:"address,omitempty"`
	JobStatus *string              `json:"job_status,omitempty"`
	Rating    float32              `json:"rating"`
}

func (r *RequestCollectorDTO) ValidateRequestColector() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.AddressId) == "" {
		errors["address_id"] = append(errors["address_id"], "address_id harus diisi")
	}
	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
