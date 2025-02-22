package model

import "time"

type Product struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null;column:id" json:"id"`
	StoreID       string         `gorm:"type:uuid;not null;column:store_id" json:"storeId"`
	Store         Store          `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ProductName   string         `gorm:"not null;column:product_name;index" json:"productName"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"productImages,omitempty"`
	Quantity      int            `gorm:"not null;column:quantity" json:"quantity"`
	Saled         int            `gorm:"default:0;column:saled" json:"saled"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp;column:created_at" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"default:current_timestamp;column:updated_at" json:"updatedAt"`
}

type ProductImage struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null;column:id" json:"id"`
	ProductID string  `gorm:"type:uuid;not null;column:product_id" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ImageURL  string  `gorm:"not null;column:image_url" json:"imageURL"`
}
