package dto

import "strings"

type RequestCoverageArea struct {
	Province string `json:"province"`
	Regency  string `json:"regency"`
}

type ResponseCoverageArea struct {
	ID        string `json:"id"`
	Province  string `json:"province"`
	Regency   string `json:"regency"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (r *RequestCoverageArea) ValidateCoverageArea() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Province) == "" {
		errors["province"] = append(errors["province"], "nama provinsi harus diisi")
	}

	if strings.TrimSpace(r.Regency) == "" {
		errors["regency"] = append(errors["regency"], "nama regency harus diisi")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
