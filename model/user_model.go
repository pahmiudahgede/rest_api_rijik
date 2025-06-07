package model

import "time"

type User struct {
	ID                   string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	Avatar               *string `json:"avatar,omitempty"`
	Name                 string  `gorm:"not null" json:"name"`
	Gender               string  `gorm:"not null" json:"gender"`
	Dateofbirth          string  `gorm:"not null" json:"dateofbirth"`
	Placeofbirth         string  `gorm:"not null" json:"placeofbirth"`
	Phone                string  `gorm:"not null;index" json:"phone"`
	Email                string  `json:"email,omitempty"`
	PhoneVerified        bool    `gorm:"default:false" json:"phoneVerified"`
	Password             string  `json:"password,omitempty"`
	RoleID               string  `gorm:"not null" json:"roleId"`
	Role                 *Role   `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`
	RegistrationStatus   string  `gorm:"default:uncompleted" json:"registrationstatus"`
	RegistrationProgress int8    `json:"registration_progress"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
