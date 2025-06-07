package article

import (
	"strings"
)

type ArticleResponseDTO struct {
	ID          string `json:"article_id"`
	Title       string `json:"title"`
	CoverImage  string `json:"coverImage"`
	Author      string `json:"author"`
	Heading     string `json:"heading"`
	Content     string `json:"content"`
	PublishedAt string `json:"publishedAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type RequestArticleDTO struct {
	Title      string `json:"title"`
	CoverImage string `json:"coverImage"`
	Author     string `json:"author"`
	Heading    string `json:"heading"`
	Content    string `json:"content"`
}

func (r *RequestArticleDTO) ValidateRequestArticleDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Title) == "" {
		errors["title"] = append(errors["title"], "Title is required")
	}

	if strings.TrimSpace(r.Author) == "" {
		errors["author"] = append(errors["author"], "Author is required")
	}
	if strings.TrimSpace(r.Heading) == "" {
		errors["heading"] = append(errors["heading"], "Heading is required")
	}
	if strings.TrimSpace(r.Content) == "" {
		errors["content"] = append(errors["content"], "Content is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
