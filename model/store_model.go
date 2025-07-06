package model

import "time"

type Store struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null;column:id" json:"id"`
	UserID         string    `gorm:"type:uuid;not null;column:user_id" json:"userId"`
	User           User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StoreName      string    `gorm:"not null;column:store_name;index" json:"storeName"`
	StoreLogo      string    `gorm:"not null;column:store_logo" json:"storeLogo"`
	StoreBanner    string    `gorm:"not null;column:store_banner" json:"storeBanner"`
	StoreInfo      string    `gorm:"not null;column:store_info" json:"storeInfo"`
	StoreAddressID string    `gorm:"type:uuid;not null;column:store_address_id" json:"storeAddressId"`
	StoreAddress   Address   `gorm:"foreignKey:StoreAddressID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Followers      int       `gorm:"default:0;column:followers" json:"followers,omitempty"`
	Products       []Product `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;" json:"products,omitempty"`
	TotalProduct   int       `gorm:"default:0;column:total_product" json:"totalProduct,omitempty"`
	CreatedAt      time.Time `gorm:"default:current_timestamp;column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp;column:updated_at" json:"updatedAt"`
}
