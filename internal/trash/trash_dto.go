package trash

import (
	"strings"
)

type RequestTrashCategoryDTO struct {
	Name           string  `json:"name"`
	EstimatedPrice float64 `json:"estimated_price"`
	IconTrash      string  `json:"icon_trash,omitempty"`
	Variety        string  `json:"variety"`
}

type RequestTrashDetailDTO struct {
	CategoryID      string `json:"category_id"`
	StepOrder       int    `json:"step"`
	IconTrashDetail string `json:"icon_trash_detail,omitempty"`
	Description     string `json:"description"`
}

type ResponseTrashCategoryDTO struct {
	ID             string                   `json:"id,omitempty"`
	TrashName      string                   `json:"trash_name,omitempty"`
	TrashIcon      string                   `json:"trash_icon,omitempty"`
	EstimatedPrice float64                  `json:"estimated_price"`
	Variety        string                   `json:"variety,omitempty"`
	CreatedAt      string                   `json:"created_at,omitempty"`
	UpdatedAt      string                   `json:"updated_at,omitempty"`
	TrashDetail    []ResponseTrashDetailDTO `json:"trash_detail,omitempty"`
}

type ResponseTrashDetailDTO struct {
	ID              string `json:"trashdetail_id"`
	CategoryID      string `json:"category_id"`
	IconTrashDetail string `json:"trashdetail_icon,omitempty"`
	Description     string `json:"description"`
	StepOrder       int    `json:"step_order"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (r *RequestTrashCategoryDTO) ValidateRequestTrashCategoryDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = append(errors["name"], "name is required")
	}
	if r.EstimatedPrice <= 0 {
		errors["estimated_price"] = append(errors["estimated_price"], "estimated price must be greater than 0")
	}
	if strings.TrimSpace(r.Variety) == "" {
		errors["variety"] = append(errors["variety"], "variety is required")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

func (r *RequestTrashDetailDTO) ValidateRequestTrashDetailDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Description) == "" {
		errors["description"] = append(errors["description"], "description is required")
	}
	if strings.TrimSpace(r.CategoryID) == "" {
		errors["category_id"] = append(errors["category_id"], "category_id is required")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
