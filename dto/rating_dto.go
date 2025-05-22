package dto

import "strings"

type CreatePickupRatingDTO struct {
	Rating   float32 `json:"rating"`
	Feedback string  `json:"feedback"`
}

func (r *CreatePickupRatingDTO) ValidateCreatePickupRatingDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.Rating < 1.0 || r.Rating > 5.0 {
		errors["rating"] = append(errors["rating"], "Rating harus antara 1.0 sampai 5.0")
	}

	if len(strings.TrimSpace(r.Feedback)) > 255 {
		errors["feedback"] = append(errors["feedback"], "Feedback tidak boleh lebih dari 255 karakter")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
