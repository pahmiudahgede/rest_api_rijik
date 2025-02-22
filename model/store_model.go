package model

import "time"

type Store struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID         string    `gorm:"type:uuid;not null" json:"userId"`
	User           User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	StoreName      string    `gorm:"not null" json:"storeName"`
	StoreLogo      string    `gorm:"not null" json:"storeLogo"`
	StoreBanner    string    `gorm:"not null" json:"storeBanner"`
	StoreInfo      string    `gorm:"not null" json:"storeInfo"`
	StoreAddressID string    `gorm:"type:uuid;not null" json:"storeAddressId"`
	StoreAddress   Address   `gorm:"foreignKey:StoreAddressID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"storeAddress"`
	Followers      int       `gorm:"default:0" json:"followers"`
	Products       []Product `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;" json:"products"`
	CreatedAt      time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
