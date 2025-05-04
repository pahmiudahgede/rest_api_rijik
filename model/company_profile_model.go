package model

import (
	"time"
)

type CompanyProfile struct {
	ID                 string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID             string    `gorm:"not null" json:"userId"`
	User               User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CompanyName        string    `gorm:"not null" json:"company_name"`
	CompanyAddress     string    `gorm:"not null" json:"company_address"`
	CompanyPhone       string    `gorm:"not null" json:"company_phone"`
	CompanyEmail       string    `gorm:"not null" json:"company_email"`
	CompanyLogo        string    `gorm:"not null" json:"company_logo"`
	CompanyWebsite     string    `json:"company_website"`
	TaxID              string    `json:"tax_id"`
	FoundedDate        time.Time `json:"founded_date"`
	CompanyType        string    `gorm:"not null" json:"company_type"`
	CompanyDescription string    `gorm:"type:text" json:"company_description"`
	CompanyStatus      string    `gorm:"not null" json:"company_status"`
	CreatedAt          time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
