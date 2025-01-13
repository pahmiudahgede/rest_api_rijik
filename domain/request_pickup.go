package domain

import "time"

type RequestPickup struct {
	ID            string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID        string        `gorm:"type:uuid;not null" json:"userId"`
	Request       []RequestItem `gorm:"foreignKey:RequestPickupID" json:"request"`
	RequestTime   string        `gorm:"type:text;not null" json:"requestTime"`
	UserAddressID string        `gorm:"type:uuid;not null" json:"userAddressId"`
	UserAddress   Address       `gorm:"foreignKey:UserAddressID" json:"userAddress"`
	StatusRequest string        `gorm:"type:text;not null" json:"statusRequest"`
	CreatedAt     time.Time     `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time     `gorm:"default:current_timestamp" json:"updatedAt"`
}

type RequestItem struct {
	ID              string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	RequestPickupID string        `gorm:"type:uuid;not null" json:"requestPickupId"`
	TrashCategoryID string        `gorm:"type:uuid;not null" json:"trashCategoryId"`
	TrashCategory   TrashCategory `gorm:"foreignKey:TrashCategoryID" json:"trashCategory"`
	EstimatedAmount string        `gorm:"type:text;not null" json:"estimatedAmount"`
}
