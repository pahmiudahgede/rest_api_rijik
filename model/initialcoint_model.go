package model

import "time"

type InitialCoint struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CoinName     string    `gorm:"not null" json:"coin_name"`
	ValuePerUnit float64   `gorm:"not null" json:"value_perunit"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
