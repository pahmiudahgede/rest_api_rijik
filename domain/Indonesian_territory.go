package domain

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Regency struct {
	ID         string    `json:"id"`
	ProvinceID string    `json:"province_id"`
	Name       string    `json:"name"`
	Province   *Province `json:"province,omitempty"`
}

type District struct {
	ID        string   `json:"id"`
	RegencyID string   `json:"regency_id"`
	Name      string   `json:"name"`
	Regency   *Regency `json:"regency,omitempty"`
}

type Village struct {
	ID         string    `json:"id"`
	DistrictID string    `json:"district_id"`
	Name       string    `json:"name"`
	District   *District `json:"district,omitempty"`
}
