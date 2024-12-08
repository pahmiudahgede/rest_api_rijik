package domain

import (
	"time"
)

type User struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	Avatar        *string   `json:"avatar,omitempty"`
	Username      string    `gorm:"unique;not null" json:"username"`
	Name      string    `gorm:"not null" json:"name"`
	Phone         string    `gorm:"not null" json:"phone"`
	Email         string    `gorm:"unique;not null" json:"email"`
	EmailVerified bool      `gorm:"default:false" json:"emailVerified"`
	Password      string    `gorm:"not null" json:"password"`
	CreatedAt     time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
	RoleID        string    `gorm:"not null" json:"roleId"`
	Role          UserRole  `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"role"`
	AddressId     *string   `gorm:"default:null" json:"addressId"`
	Addresses     []Address `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"addresses"`
}
