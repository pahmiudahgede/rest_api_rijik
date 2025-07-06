package role

type RoleResponseDTO struct {
	ID        string `json:"role_id"`
	RoleName  string `json:"role_name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
