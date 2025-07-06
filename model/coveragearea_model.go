package model

import "time"

type CoverageArea struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	Province  string    `gorm:"not null" json:"province"`
	Regency   string    `gorm:"not null" json:"regency"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
