package dto

type BannerResponse struct {
	ID          string `json:"id"`
	BannerName  string `json:"bannername"`
	BannerImage string `json:"bannerimage"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type BannerCreateRequest struct {
	BannerName  string `json:"bannername" validate:"required"`
	BannerImage string `json:"bannerimage" validate:"required"`
}

type BannerUpdateRequest struct {
	BannerName  *string `json:"bannername,omitempty"`
	BannerImage *string `json:"bannerimage,omitempty"`
}
