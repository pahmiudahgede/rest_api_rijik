package dto

import "github.com/go-playground/validator/v10"

type PointResponse struct {
	ID           string  `json:"id"`
	CoinName     string  `json:"coin_name"`
	ValuePerUnit float64 `json:"value_perunit"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func NewPointResponse(id, coinName string, valuePerUnit float64, createdAt, updatedAt string) PointResponse {
	return PointResponse{
		ID:           id,
		CoinName:     coinName,
		ValuePerUnit: valuePerUnit,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

type PointRequest struct {
	CoinName     string  `json:"coin_name" validate:"required"`
	ValuePerUnit float64 `json:"value_perunit" validate:"required,gt=0"`
}

func NewPointRequest(coinName string, valuePerUnit float64) PointRequest {
	return PointRequest{
		CoinName:     coinName,
		ValuePerUnit: valuePerUnit,
	}
}

func (p *PointRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type PointUpdateDTO struct {
	CoinName     string  `json:"coin_name" validate:"required"`
	ValuePerUnit float64 `json:"value_perunit" validate:"required,gt=0"`
}

func (p *PointUpdateDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}