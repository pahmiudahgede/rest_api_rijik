package dto

type ProductResponseDTO struct {
	ID            string                 `json:"id"`
	UserID        string                 `json:"user_id"`
	ProductTitle  string                 `json:"product_title"`
	ProductImages []ProductImageDTO      `json:"product_images"`
	TrashDetail   TrashDetailResponseDTO `json:"trash_detail"`

	SalePrice       int64  `json:"sale_price"`
	Quantity        int    `json:"quantity"`
	ProductDescribe string `json:"product_describe"`
	Sold            int    `json:"sold"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ProductImageDTO struct {
	ImageURL string `json:"image_url"`
}

type TrashDetailResponseDTO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}
