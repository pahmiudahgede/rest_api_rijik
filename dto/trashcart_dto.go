package dto

import (
	"fmt"
	"strings"
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
	CreatedAt           string             `json:"createdAt"`
	UpdatedAt           string             `json:"updatedAt"`
}

type CartItemResponse struct {
	ItemId                 string  `json:"item_id"`
	TrashId                string  `json:"trashid"`
	TrashIcon              string  `json:"trashicon"`
	TrashName              string  `json:"trashname"`
	Amount                 float32 `json:"amount"`
	EstimatedSubTotalPrice float32 `json:"estimated_subtotalprice"`
}

type RequestCartItems struct {
	TrashCategoryID string  `json:"trashid"`
	Amount          float32 `json:"amount"`
}

func (r *RequestCartItems) ValidateRequestCartItem() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.TrashCategoryID) == "" {
		errors["trashid"] = append(errors["trashid"], "trashid is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

type BulkRequestCartItems struct {
	Items []RequestCartItems `json:"items"`
}

func (b *BulkRequestCartItems) Validate() (map[string][]string, bool) {
	errors := make(map[string][]string)
	for i, item := range b.Items {
		if strings.TrimSpace(item.TrashCategoryID) == "" {
			errors[fmt.Sprintf("items[%d].trashid", i)] = append(errors[fmt.Sprintf("items[%d].trashid", i)], "trashid is required")
		}
	}
	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
