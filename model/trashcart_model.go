package model

import (
	"time"
)

type Cart struct {
	ID                  string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID              string     `gorm:"not null" json:"user_id"`
	User                User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	CartItems           []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"cart_items"`
	TotalAmount         float32    `json:"total_amount"`
	EstimatedTotalPrice float32    `json:"estimated_total_price"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type CartItem struct {
	ID                     string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID                 string        `gorm:"not null" json:"-"`
	Cart                   *Cart         `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"-"`
	TrashCategoryID        string        `gorm:"not null" json:"trash_id"`
	TrashCategory          TrashCategory `gorm:"foreignKey:TrashCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"trash_category"`
	Amount                 float32       `json:"amount"`
	SubTotalEstimatedPrice float32       `json:"subtotal_estimated_price"`
	CreatedAt              time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
