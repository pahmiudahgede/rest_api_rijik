package domain

import "time"

type TrashCategory struct {
	ID        string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string        `gorm:"not null" json:"name"`
	Details   []TrashDetail `gorm:"foreignKey:CategoryID" json:"details"`
	CreatedAt time.Time     `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time     `gorm:"default:current_timestamp" json:"updatedAt"`
}

type TrashDetail struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CategoryID  string    `gorm:"type:uuid;not null" json:"category_id"`
	Description string    `gorm:"not null" json:"description"`
	Price       int       `gorm:"not null" json:"price"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
