package dto

import "errors"

type ProductResponseWithSoldDTO struct {
	ID              string                 `json:"id"`
	StoreID         string                 `json:"store_id"`
	ProductTitle    string                 `json:"product_title"`
	ProductImages   []ProductImageDTO      `json:"product_images"`
	TrashDetail     TrashDetailResponseDTO `json:"trash_detail"`
	SalePrice       int64                  `json:"sale_price"`
	Quantity        int                    `json:"quantity"`
	ProductDescribe string                 `json:"product_describe"`
	Sold            int                    `json:"sold"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

type ProductResponseWithoutSoldDTO struct {
	ID              string                 `json:"id"`
	StoreID         string                 `json:"store_id"`
	ProductTitle    string                 `json:"product_title"`
	ProductImages   []ProductImageDTO      `json:"product_images"`
	TrashDetail     TrashDetailResponseDTO `json:"trash_detail"`
	SalePrice       int64                  `json:"sale_price"`
	Quantity        int                    `json:"quantity"`
	ProductDescribe string                 `json:"product_describe"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

type ProductResponseDTO struct {
	ID              string                 `json:"id"`
	StoreID         string                 `json:"store_id"`
	ProductTitle    string                 `json:"product_title"`
	ProductImages   []ProductImageDTO      `json:"product_images"`
	TrashDetail     TrashDetailResponseDTO `json:"trash_detail"`
	SalePrice       int64                  `json:"sale_price"`
	Quantity        int                    `json:"quantity"`
	ProductDescribe string                 `json:"product_describe"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

type ProductImageDTO struct {
	ImageURL string `json:"image_url"`
}

type TrashDetailResponseDTO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type CreateProductRequestDTO struct {
	StoreID         string   `json:"storeid" validate:"required,uuid"`
	ProductTitle    string   `json:"product_title" validate:"required,min=3,max=255"`
	ProductImages   []string `json:"product_images" validate:"required,min=1,dive,url"`
	TrashDetailID   string   `json:"trash_detail_id" validate:"required,uuid"`
	SalePrice       int64    `json:"sale_price" validate:"required,gt=0"`
	Quantity        int      `json:"quantity" validate:"required,gt=0"`
	ProductDescribe string   `json:"product_describe,omitempty"`
}

type CreateProductResponseDTO struct {
	ID              string                 `json:"id"`
	StoreID         string                 `json:"store_id"`
	ProductTitle    string                 `json:"product_title"`
	ProductImages   []string               `json:"product_images"`
	TrashDetail     TrashDetailResponseDTO `json:"trash_detail"`
	SalePrice       int64                  `json:"sale_price"`
	Quantity        int                    `json:"quantity"`
	ProductDescribe string                 `json:"product_describe,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

type UpdateProductRequestDTO struct {
	ProductTitle    string   `json:"product_title" validate:"required,min=3,max=255"`
	ProductImages   []string `json:"product_images" validate:"required,min=1,dive,url"`
	TrashDetailID   string   `json:"trash_detail_id" validate:"required,uuid"`
	SalePrice       int64    `json:"sale_price" validate:"required,gt=0"`
	Quantity        int      `json:"quantity" validate:"required,gt=0"`
	ProductDescribe string   `json:"product_describe,omitempty"`
}

func ValidateSalePrice(marketPrice, salePrice int64) error {

	if salePrice > marketPrice*2 {
		return errors.New("sale price cannot be more than twice the market price")
	}
	return nil
}
