package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type PinResponse struct {
	CreatedAt string `json:"createdAt"`
}
type PinInput struct {
	Pin string `json:"pin" validate:"required,len=6,numeric"`
}

func (p *PinInput) ValidateCreate() error {
	err := validate.Struct(p)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			case "Pin":
				return fmt.Errorf("PIN harus terdiri dari 6 digit angka")
			}
		}
	}
	return nil
}

type PinUpdateInput struct {
	OldPin string `json:"old_pin" validate:"required,len=6,numeric"`
	NewPin string `json:"new_pin" validate:"required,len=6,numeric"`
}

func (p *PinUpdateInput) ValidateUpdate() error {
	err := validate.Struct(p)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			case "OldPin":
				return fmt.Errorf("PIN lama harus terdiri dari 6 digit angka")
			case "NewPin":
				return fmt.Errorf("PIN baru harus terdiri dari 6 digit angka")
			}
		}
	}
	return nil
}
