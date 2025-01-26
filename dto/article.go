package dto

type ArticleResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	CoverImage  string `json:"coverImage"`
	Author      string `json:"author"`
	Heading     string `json:"heading"`
	Content     string `json:"content"`
	PublishedAt string `json:"publishedAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ArticleCreateRequest struct {
	Title      string `json:"title" validate:"required"`
	CoverImage string `json:"coverImage" validate:"required"`
	Author     string `json:"author" validate:"required"`
	Heading    string `json:"heading" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

type ArticleUpdateRequest struct {
	Title      string `json:"title" validate:"required"`
	CoverImage string `json:"coverImage" validate:"required"`
	Author     string `json:"author" validate:"required"`
	Heading    string `json:"heading" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

func (p *ArticleCreateRequest) Validate() error {
	validate := GetValidator()
	return validate.Struct(p)
}

func (p *ArticleUpdateRequest) Validate() error {
	validate := GetValidator()
	return validate.Struct(p)
}
