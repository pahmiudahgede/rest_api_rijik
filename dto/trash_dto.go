package dto
/* 
import (
	"strings"
)

type RequestTrashCategoryDTO struct {
	Name           string `json:"name"`
	EstimatedPrice string `json:"estimatedprice"`
	Icon           string `json:"icon"`
}

type ResponseTrashCategoryDTO struct {
	ID             string                   `json:"id,omitempty"`
	Name           string                   `json:"name,omitempty"`
	Icon           string                   `json:"icon,omitempty"`
	EstimatedPrice float64                  `json:"estimatedprice"`
	CreatedAt      string                   `json:"createdAt,omitempty"`
	UpdatedAt      string                   `json:"updatedAt,omitempty"`
	Details        []ResponseTrashDetailDTO `json:"details,omitempty"`
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

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
 */