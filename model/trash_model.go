package model

import "time"

type TrashCategory struct {
	ID             string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name           string        `gorm:"not null" json:"trash_name"`
	IconTrash      string        `json:"trash_icon,omitempty"`
	EstimatedPrice float64       `gorm:"not null" json:"estimated_price"`
	Variety        string        `gorm:"not null" json:"variety"`
	Details        []TrashDetail `gorm:"foreignKey:TrashCategoryID;constraint:OnDelete:CASCADE;" json:"trash_detail"`
	CreatedAt      time.Time     `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"default:current_timestamp" json:"updatedAt"`
}

type TrashDetail struct {
	ID              string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"trashdetail_id"`
	TrashCategoryID string    `gorm:"type:uuid;not null" json:"category_id"`
	IconTrashDetail string    `json:"trashdetail_icon,omitempty"`
	Description     string    `gorm:"not null" json:"description"`
	StepOrder       int       `gorm:"not null" json:"step_order"`
	CreatedAt       time.Time `gorm:"default:current_timestamp" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp" json:"updatedAt"`
}
