package dto

import (
	"fmt"
	"strings"
)

type RequestPickup struct {
	AddressID     string              `json:"address_id"`
	RequestMethod string              `json:"request_method"`
	EvidenceImage string              `json:"evidence_image"`
	Notes         string              `json:"notes"`
	RequestItems  []RequestPickupItem `json:"request_items"`
}

type RequestPickupItem struct {
	TrashCategoryID string  `json:"trash_category_id"`
	EstimatedAmount float64 `json:"estimated_amount"`
}

type ResponseRequestPickup struct {
	ID                     string                      `json:"id,omitempty"`
	UserId                 string                      `json:"user_id,omitempty"`
	User                   []UserResponseDTO           `json:"user,omitempty"`
	AddressID              string                      `json:"address_id,omitempty"`
	Address                []AddressResponseDTO        `json:"address,omitempty"`
	EvidenceImage          string                      `json:"evidence_image,omitempty"`
	Notes                  string                      `json:"notes,omitempty"`
	StatusPickup           string                      `json:"status_pickup,omitempty"`
	CollectorID            string                      `json:"collectorid,omitempty"`
	Collector              []ResponseCollectorDTO      `json:"collector,omitempty"`
	ConfirmedByCollectorAt string                      `json:"confirmedat,omitempty"`
	CreatedAt              string                      `json:"created_at,omitempty"`
	UpdatedAt              string                      `json:"updated_at,omitempty"`
	RequestItems           []ResponseRequestPickupItem `json:"request_items,omitempty"`
}

type ResponseRequestPickupItem struct {
	ID              string                     `json:"id,omitempty"`
	TrashCategoryID string                     `json:"trash_category_id,omitempty"`
	// TrashCategory   []ResponseTrashCategoryDTO `json:"trash_category,omitempty"`
	EstimatedAmount float64                    `json:"estimated_amount,omitempty"`
}

func (r *RequestPickup) ValidateRequestPickup() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if len(r.RequestItems) == 0 {
		errors["request_items"] = append(errors["request_items"], "At least one item must be provided")
	}

	if strings.TrimSpace(r.AddressID) == "" {
		errors["address_id"] = append(errors["address_id"], "Address ID must be provided")
	}

	for i, item := range r.RequestItems {
		itemErrors, valid := item.ValidateRequestPickupItem(i)
		if !valid {
			for field, msgs := range itemErrors {
				errors[field] = append(errors[field], msgs...)
			}
		}
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *RequestPickupItem) ValidateRequestPickupItem(index int) (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.TrashCategoryID) == "" {
		errors["trash_category_id"] = append(errors["trash_category_id"], fmt.Sprintf("Trash category ID cannot be empty (Item %d)", index+1))
	}

	if r.EstimatedAmount < 2 {
		errors["estimated_amount"] = append(errors["estimated_amount"], fmt.Sprintf("Estimated amount must be >= 2.0 kg (Item %d)", index+1))
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
