package dto

import "strings"

type RequestTrashCategoryDTO struct {
	Name string `json:"name"`
}

type ResponseTrashCategoryDTO struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	CreatedAt string                  `json:"createdAt"`
	UpdatedAt string                  `json:"updatedAt"`
	Details   []ResponseTrashDetailDTO `json:"details,omitempty"`
}

type ResponseTrashDetailDTO struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"category_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type RequestTrashDetailDTO struct {
	CategoryID  string  `json:"category_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (r *RequestTrashCategoryDTO) ValidateTrashCategoryInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = append(errors["name"], "name is required")
	}
	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *RequestTrashDetailDTO) ValidateTrashDetailInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Description) == "" {
		errors["description"] = append(errors["description"], "description is required")
	}
	if r.Price <= 0 {
		errors["price"] = append(errors["price"], "price must be greater than 0")
	}
	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
