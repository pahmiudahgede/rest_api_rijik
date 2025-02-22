package model

import "time"

type Product struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	StoreID       string         `gorm:"type:uuid;not null" json:"storeId"`
	Store         Store          `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"store"`
	ProductName   string         `gorm:"not null" json:"productName"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"productImages"`
	Quantity      int            `gorm:"not null" json:"quantity"`
	Saled         int            `gorm:"default:0" json:"saled"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"default:current_timestamp" json:"updatedAt"`
}

type ProductImage struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	ProductID string  `gorm:"type:uuid;not null" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product"`
	ImageURL  string  `gorm:"not null" json:"imageURL"`
}
