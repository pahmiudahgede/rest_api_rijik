package dto

type PointResponse struct {
	ID           string  `json:"id"`
	CoinName     string  `json:"coin_name"`
	ValuePerUnit float64 `json:"value_perunit"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type PointCreateRequest struct {
	CoinName     string  `json:"coin_name" validate:"required"`
	ValuePerUnit float64 `json:"value_perunit" validate:"required,gt=0"`
}

type PointUpdateRequest struct {
	CoinName     string  `json:"coin_name" validate:"required"`
	ValuePerUnit float64 `json:"value_perunit" validate:"required,gt=0"`
}

func (p *PointCreateRequest) Validate() error {
	validate := GetValidator()
	return validate.Struct(p)
}

func (p *PointUpdateRequest) Validate() error {
	validate := GetValidator()
	return validate.Struct(p)
}
