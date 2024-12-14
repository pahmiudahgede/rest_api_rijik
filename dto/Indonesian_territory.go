package dto

type ProvinceDetailResponse struct {
	ID             string         `json:"id"`
	ProvinsiName   string         `json:"provinsi_name"`
	ListRegency    []RegencyItem  `json:"list_regency"`
}

type RegencyDetailResponse struct {
	ID            string          `json:"id"`
	RegencyName   string          `json:"regency_name"`
	ProvinceID    string          `json:"province_id"`
	ProvinceName  string          `json:"province_name"`
	ListDistrict  []DistrictItem  `json:"list_districts"`
}

type DistrictDetailResponse struct {
	ID            string           `json:"id"`
	DistrictName  string           `json:"district_name"`
	ProvinceID    string           `json:"province_id"`
	ProvinceName  string           `json:"province_name"`
	RegencyID     string           `json:"regency_id"`
	RegencyName   string           `json:"regency_name"`
	ListVillages  []VillageItem    `json:"list_villages"`
}

type VillageDetailResponse struct {
	ID           string       `json:"id"`
	VillageName  string       `json:"village_name"`
	ProvinceID   string       `json:"province_id"`
	RegencyID    string       `json:"regency_id"`
	DistrictID   string       `json:"district_id"`
	ProvinceName string       `json:"province_name"`
	RegencyName  string       `json:"regency_name"`
	DistrictName string       `json:"district_name"`
}

type RegencyItem struct {
	ID         string `json:"id"`
	RegencyName string `json:"regency_name"`
}

type DistrictItem struct {
	ID           string `json:"id"`
	DistrictName string `json:"district_name"`
}

type VillageItem struct {
	ID          string `json:"id"`
	VillageName string `json:"village_name"`
}
