package address

import "strings"

type AddressResponseDTO struct {
	UserID     string  `json:"user_id,omitempty"`
	ID         string  `json:"address_id,omitempty"`
	Province   string  `json:"province,omitempty"`
	Regency    string  `json:"regency,omitempty"`
	District   string  `json:"district,omitempty"`
	Village    string  `json:"village,omitempty"`
	PostalCode string  `json:"postalCode,omitempty"`
	Detail     string  `json:"detail,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	CreatedAt  string  `json:"createdAt,omitempty"`
	UpdatedAt  string  `json:"updatedAt,omitempty"`
}

type CreateAddressDTO struct {
	Province   string  `json:"province_id"`
	Regency    string  `json:"regency_id"`
	District   string  `json:"district_id"`
	Village    string  `json:"village_id"`
	PostalCode string  `json:"postalCode"`
	Detail     string  `json:"detail"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
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
		errors["postalCode"] = append(errors["postalCode"], "PostalCode is required")
	} else if len(r.PostalCode) < 5 {
		errors["postalCode"] = append(errors["postalCode"], "PostalCode must be at least 5 characters")
	}

	if strings.TrimSpace(r.Detail) == "" {
		errors["detail"] = append(errors["detail"], "Detail address is required")
	}

	if r.Latitude == 0 {
		errors["latitude"] = append(errors["latitude"], "Latitude is required")
	}

	if r.Longitude == 0 {
		errors["longitude"] = append(errors["longitude"], "Longitude is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
