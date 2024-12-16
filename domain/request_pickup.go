package domain

import "time"

type RequestPickup struct {
	ID            string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID        string        `gorm:"not null" json:"userId"`
	User          User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Request       []RequestItem `gorm:"foreignKey:RequestPickupID" json:"request"`
	RequestTime   string        `json:"requestTime"`
	UserAddressID string        `json:"userAddressId"`
	UserAddress   Address       `gorm:"foreignKey:UserAddressID" json:"userAddress"`
	StatusRequest string        `json:"statusRequest"`
	CreatedAt     time.Time     `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time     `gorm:"default:current_timestamp" json:"updatedAt"`
}

type RequestItem struct {
	ID              string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TrashCategoryID string        `gorm:"type:uuid;not null" json:"trashCategoryId"`
	TrashCategory   TrashCategory `gorm:"foreignKey:TrashCategoryID" json:"trashCategory"`
	EstimatedAmount string        `gorm:"not null" json:"estimatedAmount"`
	RequestPickupID string        `gorm:"type:uuid;not null" json:"requestPickupId"`
}

