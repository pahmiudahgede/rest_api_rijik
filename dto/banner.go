package dto

import "github.com/go-playground/validator/v10"

type BannerResponse struct {
	ID          string `json:"id"`
	BannerName  string `json:"bannername"`
	BannerImage string `json:"bannerimage"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func NewBannerResponse(id, bannerName, bannerImage, createdAt, updatedAt string) BannerResponse {
	return BannerResponse{
		ID:          id,
		BannerName:  bannerName,
		BannerImage: bannerImage,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type BannerRequest struct {
	BannerName  string `json:"bannername" validate:"required"`
	BannerImage string `json:"bannerimage" validate:"required,url"`
}

func NewBannerRequest(bannerName, bannerImage string) BannerRequest {
	return BannerRequest{
		BannerName:  bannerName,
		BannerImage: bannerImage,
	}
}

func (b *BannerRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}

type BannerUpdateDTO struct {
	BannerName  string `json:"bannername" validate:"required"`
	BannerImage string `json:"bannerimage" validate:"required,url"`
}

func (b *BannerUpdateDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}