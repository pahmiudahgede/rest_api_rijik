package model

import "time"

type Cart struct {
	ID                  string      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId              string      `gorm:"not null" json:"userid"`
	User                User        `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;" json:"user"`
	CartItem            []CartItems `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;" json:"cartitems"`
	TotalAmount         float32     `json:"totalamount"`
	EstimatedTotalPrice float32     `json:"estimated_totalprice"`
}

type CartItems struct {
	ID            string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TrashId       string        `json:"trashid"`
	TrashCategory TrashCategory `gorm:"foreignKey:TrashId;constraint:OnDelete:CASCADE;" json:"trash"`
	Amount        float32       `json:"amount"`
	CreaatedAt    time.Time     `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time     `gorm:"default:current_timestamp" json:"updatedAt"`
}
