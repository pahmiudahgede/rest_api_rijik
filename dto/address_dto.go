package dto

import "strings"

type AddressResponseDTO struct {
	UserID     string `json:"user_id"`
	ID         string `json:"address_id"`
	Province   string `json:"province"`
	Regency    string `json:"regency"`
	District   string `json:"district"`
	Village    string `json:"village"`
	PostalCode string `json:"postalCode"`
	Detail     string `json:"detail"`
	Geography  string `json:"geography"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type CreateAddressDTO struct {
	Province   string `json:"province_id"`
	Regency    string `json:"regency_id"`
	District   string `json:"district_id"`
	Village    string `json:"village_id"`
	PostalCode string `json:"postalCode"`
	Detail     string `json:"detail"`
	Geography  string `json:"geography"`
}

func (r *CreateAddressDTO) ValidateAddress() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Province) == "" {
		errors["province_id"] = append(errors["province_id"], "Province ID is required")
	}
	if strings.TrimSpace(r.Regency) == "" {
		errors["regency_id"] = append(errors["regency_id"], "Regency ID is required")
	}
	if strings.TrimSpace(r.District) == "" {
		errors["district_id"] = append(errors["district_id"], "District ID is required")
	}
	if strings.TrimSpace(r.Village) == "" {
		errors["village_id"] = append(errors["village_id"], "Village ID is required")
	}
	if strings.TrimSpace(r.PostalCode) == "" {
		errors["postalCode"] = append(errors["village_id"], "PostalCode ID is required")
	} else if len(r.PostalCode) < 5 {
		errors["postalCode"] = append(errors["postalCode"], "kode pos belum sesuai")
	}
	if strings.TrimSpace(r.Detail) == "" {
		errors["detail"] = append(errors["detail"], "Detail address is required")
	}
	if strings.TrimSpace(r.Geography) == "" {
		errors["geography"] = append(errors["geography"], "Geographic coordinates are required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
