package model

import "time"

type PickupRating struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	RequestID   string    `gorm:"not null;unique" json:"request_id"`
	UserID      string    `gorm:"not null" json:"user_id"`
	CollectorID string    `gorm:"not null" json:"collector_id"`
	Rating      float32   `gorm:"not null" json:"rating"`
	Feedback    string    `json:"feedback,omitempty"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"created_at"`
}
