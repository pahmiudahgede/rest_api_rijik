package dto

type StoreResponseDTO struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	StoreName   string  `json:"store_name"`
	StoreLogo   string  `json:"store_logo"`
	StoreBanner string  `json:"store_banner"`
	StoreDesc   string  `json:"store_desc"`
	Follower    int     `json:"follower"`
	StoreRating float64 `json:"store_rating"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
