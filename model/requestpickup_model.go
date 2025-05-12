package model

import (
	"time"
)

type RequestPickup struct {
	ID                     string              `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserId                 string              `gorm:"not null" json:"user_id"`
	User                   User                `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	AddressId              string              `gorm:"not null" json:"address_id"`
	Address                Address             `gorm:"foreignKey:AddressId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address"`
	RequestItems           []RequestPickupItem `gorm:"foreignKey:RequestPickupId;constraint:OnDelete:CASCADE;" json:"request_items"`
	EvidenceImage          string              `json:"evidence_image"`
	StatusPickup           string              `gorm:"default:'waiting_collector'" json:"status_pickup"`
	CollectorID            *string             `gorm:"type:uuid" json:"collector_id,omitempty"`
	ConfirmedByCollectorAt time.Time           `gorm:"default:current_timestamp" json:"confirmed_by_collector_at,omitempty"`
	RequestMethod          string              `gorm:"not null" json:"request_method"`
	CreatedAt              time.Time           `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt              time.Time           `gorm:"default:current_timestamp" json:"updated_at"`
}

type RequestPickupItem struct {
	ID              string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	RequestPickupId string        `gorm:"not null" json:"request_pickup_id"`
	RequestPickup   RequestPickup `gorm:"foreignKey:RequestPickupId;constraint:OnDelete:CASCADE;"`
	TrashCategoryId string        `gorm:"not null" json:"trash_category_id"`
	TrashCategory   TrashCategory `gorm:"foreignKey:TrashCategoryId;constraint:OnDelete:CASCADE;"`
	EstimatedAmount float64       `gorm:"not null" json:"estimated_amount"`
}
