package domain

import "time"

type CoverageArea struct {
	ID        string           `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Province  string           `gorm:"not null" json:"province"`
	CreatedAt time.Time        `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time        `gorm:"default:current_timestamp" json:"updatedAt"`
	Details   []CoverageDetail `gorm:"foreignKey:CoverageAreaID" json:"coverage_area"`
}

type CoverageDetail struct {
	ID               string             `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CoverageAreaID   string             `gorm:"not null" json:"coverage_area_id"`
	Province         string             `gorm:"not null" json:"province"`
	District         string             `gorm:"not null" json:"district"`
	CreatedAt        time.Time          `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt        time.Time          `gorm:"default:current_timestamp" json:"updatedAt"`
	LocationSpecific []LocationSpecific `gorm:"foreignKey:CoverageDetailID" json:"location_specific"`
}

type LocationSpecific struct {
	ID               string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CoverageDetailID string    `gorm:"not null" json:"coverage_detail_id"`
	Subdistrict      string    `gorm:"not null" json:"subdistrict"`
	CreatedAt        time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
