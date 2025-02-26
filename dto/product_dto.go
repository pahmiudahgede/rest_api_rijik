package dto

import (
	"mime/multipart"
	"regexp"
	"strings"
)

type ResponseProductImageDTO struct {
	ID        string `json:"id"`
	ProductID string `json:"productId"`
	ImageURL  string `json:"imageURL"`
}

type ResponseProductDTO struct {
	ID            string                    `json:"id"`
	StoreID       string                    `json:"storeId"`
	ProductName   string                    `json:"productName"`
	Quantity      int                       `json:"quantity"`
	Saled         int                       `json:"saled"`
	ProductImages []ResponseProductImageDTO `json:"productImages,omitempty"`
	CreatedAt     string                    `json:"createdAt"`
	UpdatedAt     string                    `json:"updatedAt"`
}

type RequestProductDTO struct {
	ProductName   string                  `json:"product_name"`
	Quantity      int                     `json:"quantity"`
	ProductImages []*multipart.FileHeader `json:"product_images,omitempty"`
}

func (r *RequestProductDTO) ValidateProductInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.ProductName) == "" {
		errors["product_name"] = append(errors["product_name"], "Product name is required")
	} else if len(r.ProductName) < 3 {
		errors["product_name"] = append(errors["product_name"], "Product name must be at least 3 characters long")
	} else {
		validNameRegex := `^[a-zA-Z0-9\s_.-]+$`
		if matched, _ := regexp.MatchString(validNameRegex, r.ProductName); !matched {
			errors["product_name"] = append(errors["product_name"], "Product name can only contain letters, numbers, spaces, underscores, and dashes")
		}
	}

	if r.Quantity < 1 {
		errors["quantity"] = append(errors["quantity"], "Quantity must be at least 1")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
