package dto

import (
	"fmt"
	"strings"
)

type RequestCartItemDTO struct {
	TrashID string  `json:"trash_id"`
	Amount  float64 `json:"amount"`
}

type RequestCartDTO struct {
	CartItems []RequestCartItemDTO `json:"cart_items"`
}

type ResponseCartDTO struct {
	ID                  string                `json:"id"`
	UserID              string                `json:"user_id"`
	TotalAmount         float64               `json:"total_amount"`
	EstimatedTotalPrice float64               `json:"estimated_total_price"`
	CartItems           []ResponseCartItemDTO `json:"cart_items"`
}

type ResponseCartItemDTO struct {
	ID                     string  `json:"id"`
	TrashID                string  `json:"trash_id"`
	TrashName              string  `json:"trash_name"`
	TrashIcon              string  `json:"trash_icon"`
	TrashPrice             float64 `json:"trash_price"`
	Amount                 float64 `json:"amount"`
	SubTotalEstimatedPrice float64 `json:"subtotal_estimated_price"`
}

func (r *RequestCartDTO) ValidateRequestCartDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	for i, item := range r.CartItems {
		if strings.TrimSpace(item.TrashID) == "" {
			errors[fmt.Sprintf("cart_items[%d].trash_id", i)] = append(errors[fmt.Sprintf("cart_items[%d].trash_id", i)], "trash_id tidak boleh kosong")
		}
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
