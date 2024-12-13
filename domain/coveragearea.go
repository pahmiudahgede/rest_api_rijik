package domain

import "time"

type CoverageArea struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Province  string    `gorm:"not null" json:"province"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
	CoverageDistrics []CoverageDistric `gorm:"foreignKey:CoverageAreaID" json:"coverage_districs"`
}

type CoverageDistric struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CoverageAreaID string    `gorm:"not null;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"coverage_area_id"`
	District       string    `gorm:"not null" json:"district"`
	CreatedAt      time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
	CoverageSubdistricts []CoverageSubdistrict `gorm:"foreignKey:CoverageDistrictId" json:"subdistricts"`
}

type CoverageSubdistrict struct {
	ID                 string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CoverageAreaID     string    `gorm:"not null" json:"coverage_area_id"`
	CoverageDistrictId string    `gorm:"not null;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"coverage_district_id"`
	Subdistrict        string    `gorm:"not null" json:"subdistrict"`
	CreatedAt          time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
