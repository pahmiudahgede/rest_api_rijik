package model

import (
	"time"
)

type About struct {
	ID          string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	Title       string        `gorm:"not null" json:"title"`
	CoverImage  string        `json:"cover_image"`
	AboutDetail []AboutDetail `gorm:"foreignKey:AboutID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"about_detail"`
	CreatedAt   time.Time     `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:current_timestamp" json:"updated_at"`
}

type AboutDetail struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	AboutID     string    `gorm:"not null" json:"about_id"`
	ImageDetail string    `json:"image_detail"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
