package domain

type PlatformHandle struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	Platform    string `gorm:"not null" json:"platform"`
	Description string `gorm:"not null" json:"description"`
}
