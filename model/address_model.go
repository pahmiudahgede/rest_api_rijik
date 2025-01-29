package model

import "time"

type Address struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID      string    `gorm:"not null" json:"userId"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Province    string    `gorm:"not null" json:"province"`
	District    string    `gorm:"not null" json:"district"`
	Subdistrict string    `gorm:"not null" json:"subdistrict"`
	PostalCode  int       `gorm:"not null" json:"postalCode"`
	Village     string    `gorm:"not null" json:"village"`
	Detail      string    `gorm:"not null" json:"detail"`
	Geography   string    `gorm:"not null" json:"geography"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
