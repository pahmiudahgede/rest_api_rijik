package domain

import "time"

type Article struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	CoverImage  string    `gorm:"not null" json:"coverImage"`
	Author      string    `gorm:"not null" json:"author"`
	Heading     string    `gorm:"not null" json:"heading"`
	Content     string    `gorm:"not null" json:"content"`
	PublishedAt time.Time `gorm:"not null" json:"publishedAt"`
	UpdatedAt   time.Time `gorm:"not null" json:"updatedAt"`
}