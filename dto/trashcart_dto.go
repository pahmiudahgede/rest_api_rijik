package dto

import (
	"strings"
	"time"
)

type ValidationErrors struct {
	Errors map[string][]string
}

func (v ValidationErrors) Error() string {
	return "validation error"
}

type CartResponse struct {
	ID                  string             `json:"id"`
	UserID              string             `json:"userid"`
	CartItems           []CartItemResponse `json:"cartitems"`
	TotalAmount         float32            `json:"totalamount"`
	EstimatedTotalPrice float32            `json:"estimated_totalprice"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
}

type CartItemResponse struct {
	TrashIcon              string  `json:"trashicon"`
	TrashName              string  `json:"trashname"`
	Amount                 float32 `json:"amount"`
	EstimatedSubTotalPrice float32 `json:"estimated_subtotalprice"`
}

type RequestCartItems struct {
	TrashID string  `json:"trashid"`
	Amount  float32 `json:"amount"`
}

func (r *RequestCartItems) ValidateRequestCartItem() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.TrashID) == "" {
		errors["trashid"] = append(errors["trashid"], "trashid is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
