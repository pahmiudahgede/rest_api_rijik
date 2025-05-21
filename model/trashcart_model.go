package model

import (
	"time"
)

type Cart struct {
	ID                  string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID              string     `gorm:"not null" json:"userid"`
	User                User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	CartItems           []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"cartitems"`
	TotalAmount         float32    `json:"totalamount"`
	EstimatedTotalPrice float32    `json:"estimated_totalprice"`
	CreatedAt           time.Time  `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"default:current_timestamp" json:"updated_at"`
}

type CartItem struct {
	ID                     string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID                 string        `gorm:"not null" json:"-"`
	Cart                   *Cart         `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"-"`
	TrashCategoryID        string        `gorm:"not null" json:"trash_id"`
	TrashCategory          TrashCategory `gorm:"foreignKey:TrashCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"trash_category"`
	Amount                 float32       `json:"amount"`
	SubTotalEstimatedPrice float32       `json:"subtotalestimatedprice"`
	CreatedAt              time.Time     `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt              time.Time     `gorm:"default:current_timestamp" json:"updated_at"`
}
