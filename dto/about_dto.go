package dto

import (
	"strings"
)

type RequestAboutDTO struct {
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	// AboutDetail []RequestAboutDetailDTO `json:"about_detail"`
}

func (r *RequestAboutDTO) ValidateAbout() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Title) == "" {
		errors["title"] = append(errors["title"], "Title is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

type ResponseAboutDTO struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	CoverImage  string                    `json:"cover_image"`
	AboutDetail *[]ResponseAboutDetailDTO `json:"about_detail"`
	CreatedAt   string                    `json:"created_at"`
	UpdatedAt   string                    `json:"updated_at"`
}

type RequestAboutDetailDTO struct {
	AboutId     string `json:"about_id"`
	ImageDetail string `json:"image_detail"`
	Description string `json:"description"`
}

func (r *RequestAboutDetailDTO) ValidateAboutDetail() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.AboutId) == "" {
		errors["about_id"] = append(errors["about_id"], "About ID is required")
	}

	if strings.TrimSpace(r.ImageDetail) == "" {
		errors["image_detail"] = append(errors["image_detail"], "Image detail is required")
	}

	if strings.TrimSpace(r.Description) == "" {
		errors["description"] = append(errors["description"], "Description is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

type ResponseAboutDetailDTO struct {
	ID          string `json:"id"`
	AboutID     string `json:"about_id"`
	ImageDetail string `json:"image_detail"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
