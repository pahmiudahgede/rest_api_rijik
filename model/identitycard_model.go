package model

import "time"

type IdentityCard struct {
	ID                  string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID              string    `gorm:"not null" json:"userId"`
	User                User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Identificationumber string    `gorm:"not null" json:"identificationumber"`
	Placeofbirth        string    `gorm:"not null" json:"placeofbirth"`
	Dateofbirth         string    `gorm:"not null" json:"dateofbirth"`
	Gender              string    `gorm:"not null" json:"gender"`
	BloodType           string    `gorm:"not null" json:"bloodtype"`
	District            string    `gorm:"not null" json:"district"`
	Village             string    `gorm:"not null" json:"village"`
	Neighbourhood       string    `gorm:"not null" json:"neighbourhood"`
	Religion            string    `gorm:"not null" json:"religion"`
	Maritalstatus       string    `gorm:"not null" json:"maritalstatus"`
	Job                 string    `gorm:"not null" json:"job"`
	Citizenship         string    `gorm:"not null" json:"citizenship"`
	Validuntil          string    `gorm:"not null" json:"validuntil"`
	Cardphoto           string    `gorm:"not null" json:"cardphoto"`
	CreatedAt           time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
