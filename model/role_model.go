package model

import "time"

type Role struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	RoleName  string    `gorm:"unique;not null" json:"roleName"`
	Users     []User    `gorm:"foreignKey:RoleID" json:"users"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
