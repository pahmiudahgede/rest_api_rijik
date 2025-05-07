package model

import (
	"time"
)

type RequestPickup struct {
	ID            string              `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserId        string              `gorm:"not null" json:"user_id"`
	AddressId     string              `gorm:"not null" json:"address_id"`
	RequestItems  []RequestPickupItem `gorm:"foreignKey:RequestPickupId;constraint:OnDelete:CASCADE;" json:"request_items"`
	EvidenceImage string              `json:"evidence_image"`
	StatusPickup  string              `gorm:"default:'waiting_pengepul'" json:"status_pickup"`
	CreatedAt     time.Time           `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time           `gorm:"default:current_timestamp" json:"updated_at"`
}

type RequestPickupItem struct {
	ID              string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	RequestPickupId string  `gorm:"not null" json:"request_pickup_id"`
	TrashCategoryId string  `gorm:"not null" json:"trash_category_id"`
	TrashDetailId   string  `json:"trash_detail_id,omitempty"`
	EstimatedAmount float64 `gorm:"not null" json:"estimated_amount"`
}
