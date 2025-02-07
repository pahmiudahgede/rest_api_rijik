package model

import "time"

type Address struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID     string    `gorm:"not null" json:"userId"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Province   string    `gorm:"not null" json:"province"`
	Regency    string    `gorm:"not null" json:"regency"`
	District   string    `gorm:"not null" json:"district"`
	Village    string    `gorm:"not null" json:"village"`
	PostalCode string    `gorm:"not null" json:"postalCode"`
	Detail     string    `gorm:"not null" json:"detail"`
	Geography  string    `gorm:"not null" json:"geography"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
