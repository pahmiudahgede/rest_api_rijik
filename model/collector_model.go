package model

import "time"

type Collector struct {
	ID                      string                    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID                  string                    `gorm:"not null" json:"user_id"`
	User                    User                      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	JobStatus               string                    `gorm:"default:inactive" json:"jobstatus"`
	Rating                  float32                   `gorm:"default:5" json:"rating"`
	AddressID               string                    `gorm:"not null" json:"address_id"`
	Address                 Address                   `gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address"`
	AvaibleTrashByCollector []AvaibleTrashByCollector `gorm:"foreignKey:CollectorID;constraint:OnDelete:CASCADE;" json:"avaible_trash"`
	CreatedAt               time.Time                 `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt               time.Time                 `gorm:"default:current_timestamp" json:"updated_at"`
}

type AvaibleTrashByCollector struct {
	ID              string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	CollectorID     string        `gorm:"not null" json:"collector_id"`
	Collector       *Collector    `gorm:"foreignKey:CollectorID;constraint:OnDelete:CASCADE;" json:"-"`
	TrashCategoryID string        `gorm:"not null" json:"trash_id"`
	TrashCategory   TrashCategory `gorm:"foreignKey:TrashCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"trash_category"`
	Price           float32       `json:"price"`
}
