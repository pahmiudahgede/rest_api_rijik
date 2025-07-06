package role

import (
	"context"
	"fmt"
	"time"

	"rijig/utils"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]RoleResponseDTO, error)
	GetRoleByID(ctx context.Context, roleID string) (*RoleResponseDTO, error)
}

type roleService struct {
	RoleRepo RoleRepository
}

func NewRoleService(roleRepo RoleRepository) RoleService {
	return &roleService{roleRepo}
}

func (s *roleService) GetRoles(ctx context.Context) ([]RoleResponseDTO, error) {
	cacheKey := "roles_list"

	var cachedRoles []RoleResponseDTO
	err := utils.GetCache(cacheKey, &cachedRoles)
	if err == nil {
		return cachedRoles, nil
	}

	roles, err := s.RoleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %v", err)
	}

	var roleDTOs []RoleResponseDTO
	for _, role := range roles {
		createdAt, _ := utils.FormatDateToIndonesianFormat(role.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(role.UpdatedAt)

		roleDTOs = append(roleDTOs, RoleResponseDTO{
			ID:        role.ID,
			RoleName:  role.RoleName,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	err = utils.SetCache(cacheKey, roleDTOs, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching roles data to Redis: %v\n", err)
	}

	return roleDTOs, nil
}

func (s *roleService) GetRoleByID(ctx context.Context, roleID string) (*RoleResponseDTO, error) {
	cacheKey := fmt.Sprintf("role:%s", roleID)

	var cachedRole RoleResponseDTO
	err := utils.GetCache(cacheKey, &cachedRole)
	if err == nil {
		return &cachedRole, nil
	}

	role, err := s.RoleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(role.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(role.UpdatedAt)

	roleDTO := &RoleResponseDTO{
		ID:        role.ID,
		RoleName:  role.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err = utils.SetCache(cacheKey, roleDTO, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching role data to Redis: %v\n", err)
	}

	return roleDTO, nil
}
