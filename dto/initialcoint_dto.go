package dto

import "strings"

type ReponseInitialCointDTO struct {
	ID           string  `json:"coin_id"`
	CoinName     string  `json:"coin_name"`
	ValuePerUnit float64 `json:"value_perunit"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type RequestInitialCointDTO struct {
	CoinName     string  `json:"coin_name"`
	ValuePerUnit float64 `json:"value_perunit"`
}

func (r *RequestInitialCointDTO) ValidateCointInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.CoinName) == "" {
		errors["coin_name"] = append(errors["coin_name"], "nama coin harus diisi")
	}

	if r.ValuePerUnit <= 0 {
		errors["value_perunit"] = append(errors["value_perunit"], "value per unit harus lebih besar dari 0")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
