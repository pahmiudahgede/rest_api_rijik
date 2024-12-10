package dto

import "github.com/go-playground/validator/v10"

type TrashCategoryDTO struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

func (t *TrashCategoryDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

type TrashDetailDTO struct {
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Description string `json:"description" validate:"required,min=3,max=255"`
	Price       int    `json:"price" validate:"required,min=0"`
}

func (t *TrashDetailDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

type TrashCategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewTrashCategoryResponse(id, name, createdAt, updatedAt string) TrashCategoryResponse {
	return TrashCategoryResponse{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type TrashDetailResponse struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func NewTrashDetailResponse(id, description string, price int, createdAt, updatedAt string) TrashDetailResponse {
	return TrashDetailResponse{
		ID:          id,
		Description: description,
		Price:       price,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
