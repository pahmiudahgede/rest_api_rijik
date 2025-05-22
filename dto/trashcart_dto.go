package dto

import (
	"fmt"
	"strings"
)

type RequestCartItemDTO struct {
	TrashID string  `json:"trash_id"`
	Amount  float32 `json:"amount"`
}

type RequestCartDTO struct {
	CartItems []RequestCartItemDTO `json:"cart_items"`
}

type ResponseCartItemDTO struct {
	ID                     string  `json:"id"`
	TrashID                string  `json:"trash_id"`
	TrashName              string  `json:"trash_name"`
	TrashIcon              string  `json:"trash_icon"`
	Amount                 float32 `json:"amount"`
	SubTotalEstimatedPrice float32 `json:"subtotal_estimated_price"`
}

type ResponseCartDTO struct {
	ID                  string                `json:"id"`
	UserID              string                `json:"user_id"`
	TotalAmount         float32               `json:"total_amount"`
	EstimatedTotalPrice float32               `json:"estimated_total_price"`
	CartItems           []ResponseCartItemDTO `json:"cart_items"`
}

func (r *RequestCartDTO) ValidateRequestCartDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if len(r.CartItems) == 0 {
		errors["cart_items"] = append(errors["cart_items"], "minimal satu item harus dimasukkan")
	}

	for i, item := range r.CartItems {
		if strings.TrimSpace(item.TrashID) == "" {
			errors[fmt.Sprintf("cart_items[%d].trash_id", i)] = append(errors[fmt.Sprintf("cart_items[%d].trash_id", i)], "trash_id tidak boleh kosong")
		}
		if item.Amount <= 0 {
			errors[fmt.Sprintf("cart_items[%d].amount", i)] = append(errors[fmt.Sprintf("cart_items[%d].amount", i)], "amount harus lebih dari 0")
		}
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
