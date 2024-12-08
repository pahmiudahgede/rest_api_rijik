package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AddressInput struct {
	Province    string `json:"province" validate:"required"`
	District    string `json:"district" validate:"required"`
	Subdistrict string `json:"subdistrict" validate:"required"`
	PostalCode  int    `json:"postalCode" validate:"required,numeric"`
	Village     string `json:"village" validate:"required"`
	Detail      string `json:"detail" validate:"required"`
	Geography   string `json:"geography" validate:"required"`
}

var validate = validator.New()

func (c *AddressInput) ValidatePost() error {
	err := validate.Struct(c)
	if err != nil {

		for _, e := range err.(validator.ValidationErrors) {

			switch e.Field() {
			case "Province":
				return fmt.Errorf("provinsi harus diisisi")
			case "District":
				return fmt.Errorf("kabupaten harus diisi")
			case "Subdistrict":
				return fmt.Errorf("kecamatan harus diisi")
			case "PostalCode":
				return fmt.Errorf("postal code harus diisi dan berupa angka")
			case "Village":
				return fmt.Errorf("desa harus diisi")
			case "Detail":
				return fmt.Errorf("detail wajib diisi")
			case "Geography":
				return fmt.Errorf("lokasi kordinat harus diisi")
			}
		}
	}
	return nil
}

func (c *AddressInput) ValidateUpdate() error {
	err := validate.Struct(c)
	if err != nil {

		for _, e := range err.(validator.ValidationErrors) {

			switch e.Field() {
			case "Province":
				return fmt.Errorf("provinsi harus diisisi")
			case "District":
				return fmt.Errorf("kabupaten harus diisi")
			case "Subdistrict":
				return fmt.Errorf("kecamatan harus diisi")
			case "PostalCode":
				return fmt.Errorf("postal code harus diisi dan berupa angka")
			case "Village":
				return fmt.Errorf("desa harus diisi")
			case "Detail":
				return fmt.Errorf("detail wajib diisi")
			case "Geography":
				return fmt.Errorf("lokasi kordinat harus diisi")
			}
		}
	}
	return nil
}
