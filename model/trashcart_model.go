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
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

type CartItem struct {
	ID                     string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID                 string        `gorm:"not null" json:"-"`
	TrashID                string        `gorm:"not null" json:"trashid"`
	TrashCategory          TrashCategory `gorm:"foreignKey:TrashID;constraint:OnDelete:CASCADE;" json:"trash"`
	Amount                 float32       `json:"amount"`
	SubTotalEstimatedPrice float32       `json:"subtotalestimatedprice"`
	CreatedAt              time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`
}
