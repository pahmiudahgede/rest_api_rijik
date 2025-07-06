package requestpickup

import (
	"strings"
)

type SelectCollectorDTO struct {
	CollectorID string `json:"collector_id"`
}

type UpdateRequestPickupItemDTO struct {
	ItemID string  `json:"item_id"`
	Amount float64 `json:"actual_amount"`
}

type UpdatePickupItemsRequest struct {
	Items []UpdateRequestPickupItemDTO `json:"items"`
}

func (r *SelectCollectorDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.CollectorID) == "" {
		errors["collector_id"] = append(errors["collector_id"], "collector_id tidak boleh kosong")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

type AssignedPickupDTO struct {
	PickupID     string   `json:"pickup_id"`
	UserID       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Notes        string   `json:"notes"`
	MatchedTrash []string `json:"matched_trash"`
}

type PickupRequestForCollectorDTO struct {
	PickupID     string   `json:"pickup_id"`
	UserID       string   `json:"user_id"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	DistanceKm   float64  `json:"distance_km"`
	MatchedTrash []string `json:"matched_trash"`
}

type RequestPickupDTO struct {
	AddressID     string `json:"address_id"`
	RequestMethod string `json:"request_method"`
	Notes         string `json:"notes,omitempty"`
}

func (r *RequestPickupDTO) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.AddressID) == "" {
		errors["address_id"] = append(errors["address_id"], "alamat harus dipilih")
	}

	method := strings.ToLower(strings.TrimSpace(r.RequestMethod))
	if method != "manual" && method != "otomatis" {
		errors["request_method"] = append(errors["request_method"], "harus manual atau otomatis")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
