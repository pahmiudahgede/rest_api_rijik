package domain

import "time"

type MenuAccess struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	RoleID    string    `gorm:"not null" json:"roleId"`
	Role      UserRole  `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"role"`
	MenuName  string    `gorm:"not null" json:"menuName"`
	Path      string    `gorm:"not null" json:"path"`
	IconURL   string    `gorm:"not null" json:"iconUrl"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
