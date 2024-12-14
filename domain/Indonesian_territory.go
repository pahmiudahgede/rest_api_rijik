package domain

type Province struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ListRegency []Regency `json:"list_regency,omitempty"`
}

type Regency struct {
	ID           string     `json:"id"`
	ProvinceID   string     `json:"province_id"`
	Name         string     `json:"name"`
	Province     *Province  `json:"province,omitempty"`
	ListDistrict []District `json:"list_district,omitempty"`
}

type District struct {
	ID          string    `json:"id"`
	RegencyID   string    `json:"regency_id"`
	Name        string    `json:"name"`
	Regency     *Regency  `json:"regency,omitempty"`
	ListVillage []Village `json:"list_village,omitempty"`
}

type Village struct {
	ID         string    `json:"id"`
	DistrictID string    `json:"district_id"`
	Name       string    `json:"name"`
	District   *District `json:"district,omitempty"`
}
