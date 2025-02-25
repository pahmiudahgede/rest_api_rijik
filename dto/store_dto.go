package dto

import (
	"regexp"
	"strings"
)

type ResponseStoreDTO struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	StoreName      string `json:"storeName"`
	StoreLogo      string `json:"storeLogo"`
	StoreBanner    string `json:"storeBanner"`
	StoreInfo      string `json:"storeInfo"`
	StoreAddressID string `json:"storeAddressId"`
	TotalProduct   int    `json:"TotalProduct"`
	Followers      int    `json:"followers"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type RequestStoreDTO struct {
	StoreName      string `json:"store_name"`
	StoreLogo      string `json:"store_logo"`
	StoreBanner    string `json:"store_banner"`
	StoreInfo      string `json:"store_info"`
	StoreAddressID string `json:"store_address_id"`
}

func (r *RequestStoreDTO) ValidateStoreInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.StoreName) == "" {
		errors["store_name"] = append(errors["store_name"], "Store name is required")
	} else if len(r.StoreName) < 3 {
		errors["store_name"] = append(errors["store_name"], "Store name must be at least 3 characters long")
	} else {
		validNameRegex := `^[a-zA-Z0-9_.\s]+$`
		if matched, _ := regexp.MatchString(validNameRegex, r.StoreName); !matched {
			errors["store_name"] = append(errors["store_name"], "Store name can only contain letters, numbers, underscores, and periods")
		}
	}

	if strings.TrimSpace(r.StoreLogo) == "" {
		errors["store_logo"] = append(errors["store_logo"], "Store logo is required")
	}

	if strings.TrimSpace(r.StoreBanner) == "" {
		errors["store_banner"] = append(errors["store_banner"], "Store banner is required")
	}

	if strings.TrimSpace(r.StoreInfo) == "" {
		errors["store_info"] = append(errors["store_info"], "Store info is required")
	}

	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`
	if r.StoreAddressID == "" {
		errors["store_address_id"] = append(errors["store_address_id"], "Store address ID is required")
	} else if matched, _ := regexp.MatchString(uuidRegex, r.StoreAddressID); !matched {
		errors["store_address_id"] = append(errors["store_address_id"], "Invalid Store Address ID format")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
