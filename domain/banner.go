package domain

import "time"

type Banner struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	BannerName  string    `gorm:"not null" json:"bannername"`
	BannerImage string    `gorm:"not null" json:"bannerimage"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
