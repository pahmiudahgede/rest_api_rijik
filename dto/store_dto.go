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
		errors["storeName"] = append(errors["storeName"], "Store name is required")
	} else if len(r.StoreName) < 3 {
		errors["storeName"] = append(errors["storeName"], "Store name must be at least 3 characters long")
	} else {
		validNameRegex := `^[a-zA-Z0-9_.\s]+$`
		if matched, _ := regexp.MatchString(validNameRegex, r.StoreName); !matched {
			errors["storeName"] = append(errors["storeName"], "Store name can only contain letters, numbers, underscores, and periods")
		}
	}

	if strings.TrimSpace(r.StoreLogo) == "" {
		errors["storeLogo"] = append(errors["storeLogo"], "Store logo is required")
	}

	if strings.TrimSpace(r.StoreBanner) == "" {
		errors["storeBanner"] = append(errors["storeBanner"], "Store banner is required")
	}

	if strings.TrimSpace(r.StoreInfo) == "" {
		errors["storeInfo"] = append(errors["storeInfo"], "Store info is required")
	}

	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`
	if r.StoreAddressID == "" {
		errors["storeAddressId"] = append(errors["storeAddressId"], "Store address ID is required")
	} else if matched, _ := regexp.MatchString(uuidRegex, r.StoreAddressID); !matched {
		errors["storeAddressId"] = append(errors["storeAddressId"], "Invalid Store Address ID format")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
