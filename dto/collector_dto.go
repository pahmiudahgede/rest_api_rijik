package dto

import "strings"

type RequestCollectorDTO struct {
	UserId    string `json:"user_id"`
	AddressId string `json:"address_id"`
}

type ResponseCollectorDTO struct {
	ID        string  `json:"collector_id"`
	UserId    string  `json:"user_id"`
	AddressId string  `json:"address_id"`
	JobStatus string  `json:"job_status"`
	Rating    float32 `json:"rating"`
	// CreatedAt string  `json:"createdAt"`
	// UpdatedAt string  `json:"updatedAt"`
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
