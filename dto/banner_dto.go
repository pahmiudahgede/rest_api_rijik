package dto

import "strings"

type ResponseBannerDTO struct {
	ID          string `json:"id"`
	BannerName  string `json:"bannername"`
	BannerImage string `json:"bannerimage"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type RequestBannerDTO struct {
	BannerName  string `json:"bannername"`
	BannerImage string `json:"bannerimage"`
}

func (r *RequestBannerDTO) ValidateBannerInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.BannerName) == "" {
		errors["bannername"] = append(errors["bannername"], "nama banner harus diisi")
	}
	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
