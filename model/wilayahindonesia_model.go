package model

type Province struct {
	ID   string `gorm:"primaryKey;type:varchar(255);not null" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	Regencies []Regency `gorm:"foreignKey:ProvinceID" json:"regencies"`
}

type Regency struct {
	ID         string `gorm:"primaryKey;type:varchar(255);not null" json:"id"`
	ProvinceID string `gorm:"type:varchar(255);not null" json:"province_id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`

	Province  Province   `gorm:"foreignKey:ProvinceID" json:"province"`
	Districts []District `gorm:"foreignKey:RegencyID" json:"districts"`
}

type District struct {
	ID        string `gorm:"primaryKey;type:varchar(255);not null" json:"id"`
	RegencyID string `gorm:"type:varchar(255);not null" json:"regency_id"`
	Name      string `gorm:"type:varchar(255);not null" json:"name"`

	Regency  Regency   `gorm:"foreignKey:RegencyID" json:"regency"`
	Villages []Village `gorm:"foreignKey:DistrictID" json:"villages"`
}

type Village struct {
	ID         string `gorm:"primaryKey;type:varchar(255);not null" json:"id"`
	DistrictID string `gorm:"type:varchar(255);not null" json:"district_id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`

	District District `gorm:"foreignKey:DistrictID" json:"district"`
}
