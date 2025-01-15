package domain

import "time"

type Product struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID          string         `gorm:"type:uuid;not null" json:"user_id"`
	ProductTitle    string         `gorm:"not null" json:"product_title"`
	ProductImages   []ProductImage `gorm:"foreignKey:ProductID" json:"product_images"`
	TrashDetailID   string         `gorm:"type:uuid;not null" json:"trash_detail_id"`
	TrashDetail     TrashDetail    `gorm:"foreignKey:TrashDetailID" json:"trash_detail"`
	SalePrice       int64          `gorm:"not null" json:"sale_price"`
	Quantity        int            `gorm:"not null" json:"quantity"`
	ProductDescribe string         `gorm:"type:text" json:"product_describe"`
	Sold            int            `gorm:"default:0" json:"sold"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
}

type ProductImage struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProductID string    `gorm:"type:uuid;not null" json:"product_id"`
	ImageURL  string    `gorm:"not null" json:"image_url"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
