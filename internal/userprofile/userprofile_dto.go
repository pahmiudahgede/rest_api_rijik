package userprofile

import (
	"rijig/internal/role"
	"rijig/utils"
	"strings"
)

type UserProfileResponseDTO struct {
	ID            string               `json:"id,omitempty"`
	Avatar        string               `json:"avatar,omitempty"`
	Name          string               `json:"name,omitempty"`
	Gender        string               `json:"gender,omitempty"`
	Dateofbirth   string               `json:"dateofbirth,omitempty"`
	Placeofbirth  string               `json:"placeofbirth,omitempty"`
	Phone         string               `json:"phone,omitempty"`
	Email         string               `json:"email,omitempty"`
	PhoneVerified bool                 `json:"phone_verified,omitempty"`
	Password      string               `json:"password,omitempty"`
	Role          role.RoleResponseDTO `json:"role"`
	CreatedAt     string               `json:"createdAt,omitempty"`
	UpdatedAt     string               `json:"updatedAt,omitempty"`
}

type RequestUserProfileDTO struct {
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	Dateofbirth  string `json:"dateofbirth"`
	Placeofbirth string `json:"placeofbirth"`
	Phone        string `json:"phone"`
}

func (r *RequestUserProfileDTO) ValidateRequestUserProfileDTO() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	}

	if strings.TrimSpace(r.Gender) == "" {
		errors["gender"] = append(errors["gender"], "jenis kelamin tidak boleh kosong")
	} else if r.Gender != "perempuan" && r.Gender != "laki-laki" {
		errors["gender"] = append(errors["gender"], "jenis kelamin harus 'perempuan' atau 'laki-laki'")
	}

	if strings.TrimSpace(r.Dateofbirth) == "" {
		errors["dateofbirth"] = append(errors["dateofbirth"], "tanggal lahir dibutuhkan")
	} else if !utils.IsValidDate(r.Dateofbirth) {
		errors["dateofbirth"] = append(errors["dateofbirth"], "tanggal lahir harus berformat DD-MM-YYYY")
	}

	if strings.TrimSpace(r.Placeofbirth) == "" {
		errors["placeofbirth"] = append(errors["placeofbirth"], "Name is required")
	}

	if strings.TrimSpace(r.Phone) == "" {
		errors["phone"] = append(errors["phone"], "Phone number is required")
	} else if !utils.IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Invalid phone number format. Use 62 followed by 9-13 digits")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
