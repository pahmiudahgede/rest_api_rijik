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
	CompanyEmail       string    `json:"company_email,omitempty"`
	CompanyLogo        string    `json:"company_logo,omitempty"`
	CompanyWebsite     string    `json:"company_website,omitempty"`
	TaxID              string    `json:"tax_id,omitempty"`
	FoundedDate        string    `json:"founded_date,omitempty"`
	CompanyType        string    `json:"company_type,omitempty"`
	CompanyDescription string    `gorm:"type:text" json:"company_description"`
	CreatedAt          time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
