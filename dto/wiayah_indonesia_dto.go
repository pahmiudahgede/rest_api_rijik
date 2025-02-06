package dto

type ProvinceResponseDTO struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Regencies []RegencyResponseDTO `json:"regencies,omitempty"`
}

type RegencyResponseDTO struct {
	ID         string                `json:"id"`
	ProvinceID string                `json:"province_id"`
	Name       string                `json:"name"`
	Districts  []DistrictResponseDTO `json:"districts,omitempty"`
}

type DistrictResponseDTO struct {
	ID        string               `json:"id"`
	RegencyID string               `json:"regency_id"`
	Name      string               `json:"name"`
	Villages  []VillageResponseDTO `json:"villages,omitempty"`
}

type VillageResponseDTO struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
}
