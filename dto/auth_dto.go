package dto

type LoginDTO struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
}

type UserResponseWithToken struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}
