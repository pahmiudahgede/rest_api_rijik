package services

import (
	"context"
	"fmt"
	"time"

	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/utils"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]dto.RoleResponseDTO, error)
	GetRoleByID(ctx context.Context, roleID string) (*dto.RoleResponseDTO, error)
}

type roleService struct {
	RoleRepo repositories.RoleRepository
}

func NewRoleService(roleRepo repositories.RoleRepository) RoleService {
	return &roleService{RoleRepo: roleRepo}
}

func (s *roleService) GetRoles(ctx context.Context) ([]dto.RoleResponseDTO, error) {
	cacheKey := "roles_list"
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var roles []dto.RoleResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				role, ok := item.(map[string]interface{})
				if ok {
					roles = append(roles, dto.RoleResponseDTO{
						ID:        role["role_id"].(string),
						RoleName:  role["role_name"].(string),
						CreatedAt: role["createdAt"].(string),
						UpdatedAt: role["updatedAt"].(string),
					})
				}
			}
			return roles, nil
		}
	}

	roles, err := s.RoleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %v", err)
	}

	var roleDTOs []dto.RoleResponseDTO
	for _, role := range roles {
		createdAt, _ := utils.FormatDateToIndonesianFormat(role.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(role.UpdatedAt)

		roleDTOs = append(roleDTOs, dto.RoleResponseDTO{
			ID:        role.ID,
			RoleName:  role.RoleName,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	cacheData := map[string]interface{}{
		"data": roleDTOs,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching roles data to Redis: %v\n", err)
	}

	return roleDTOs, nil
}

func (s *roleService) GetRoleByID(ctx context.Context, roleID string) (*dto.RoleResponseDTO, error) {

	role, err := s.RoleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(role.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(role.UpdatedAt)

	roleDTO := &dto.RoleResponseDTO{
		ID:        role.ID,
		RoleName:  role.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	cacheKey := fmt.Sprintf("role:%s", roleID)
	cacheData := map[string]interface{}{
		"data": roleDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching role data to Redis: %v\n", err)
	}

	return roleDTO, nil
}
