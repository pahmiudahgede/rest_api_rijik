package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

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
	PublishedAtFormatted string    `json:"publishedAtt"`
	UpdatedAtFormatted   string    `json:"updatedAtt"`
}

type FormattedResponse struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	CoverImage           string `json:"coverImage"`
	Author               string `json:"author"`
	Heading              string `json:"heading"`
	Content              string `json:"content"`
	PublishedAtFormatted string `json:"publishedAt"`
	UpdatedAtFormatted   string `json:"updatedAt"`
}
type ArticleUpdateRequest struct {
	Title      string `json:"title" validate:"required"`
	CoverImage string `json:"coverImage" validate:"required"`
	Author     string `json:"author" validate:"required"`
	Heading    string `json:"heading" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

func (c *ArticleRequest) ValidatePostArticle() error {
	err := validate.Struct(c)
	if err != nil {

		for _, e := range err.(validator.ValidationErrors) {

			switch e.Field() {
			case "Title":
				return fmt.Errorf("judul harus diisi")
			case "CoverImage":
				return fmt.Errorf("gambar cover harus diisi")
			case "Author":
				return fmt.Errorf("penulis harus diisi")
			case "Heading":
				return fmt.Errorf("heading harus diisi")
			case "Content":
				return fmt.Errorf("konten artikel harus diisi")
			}
		}
	}
	return nil
}

func (c *ArticleUpdateRequest) ValidateUpdateArticle() error {
	err := validate.Struct(c)
	if err != nil {

		for _, e := range err.(validator.ValidationErrors) {

			switch e.Field() {
			case "Title":
				return fmt.Errorf("judul harus diisi")
			case "CoverImage":
				return fmt.Errorf("gambar cover harus diisi")
			case "Author":
				return fmt.Errorf("penulis harus diisi")
			case "Heading":
				return fmt.Errorf("heading harus diisi")
			case "Content":
				return fmt.Errorf("konten artikel harus diisi")
			}
		}
	}
	return nil
}
