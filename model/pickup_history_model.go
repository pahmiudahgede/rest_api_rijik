package model

import "time"

type PickupStatusHistory struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	RequestID     string    `gorm:"not null" json:"request_id"`
	Status        string    `gorm:"not null" json:"status"`
	ChangedAt     time.Time `gorm:"not null" json:"changed_at"`
	ChangedByID   string    `gorm:"not null" json:"changed_by_id"`
	ChangedByRole string    `gorm:"not null" json:"changed_by_role"`
}
