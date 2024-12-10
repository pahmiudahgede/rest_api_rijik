package dto

import "time"

type ArticleRequest struct {
	Title      string `json:"title" validate:"required"`
	CoverImage string `json:"coverImage" validate:"required"`
	Author     string `json:"author" validate:"required"`
	Heading    string `json:"heading" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

type ArticleResponse struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	CoverImage           string    `json:"coverImage"`
	Author               string    `json:"author"`
	Heading              string    `json:"heading"`
	Content              string    `json:"content"`
	PublishedAt          time.Time `json:"publishedAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	PublishedAtFormatted string    `json:"publishedAtFormatted"`
	UpdatedAtFormatted   string    `json:"updatedAtFormatted"`
}

type FormattedResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"coverImage"`
	Author     string `json:"author"`
	Heading    string `json:"heading"`
	Content    string `json:"content"`

	PublishedAtFormatted string `json:"publishedAtFormatted"`
	UpdatedAtFormatted   string `json:"updatedAtFormatted"`
}
type ArticleUpdateRequest struct {
	Title      string `json:"title" validate:"required"`
	CoverImage string `json:"coverImage" validate:"required"`
	Author     string `json:"author" validate:"required"`
	Heading    string `json:"heading" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

func (ar *ArticleRequest) Validate() error {
	return validate.Struct(ar)
}

func FormatDateToIndonesianFormat(t time.Time) string {
	return t.Format("02-01-2006 15:04")
}
