package role

import (
	"context"
	"rijig/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(ctx context.Context, id string) (*model.Role, error)
	FindRoleByName(ctx context.Context, roleName string) (*model.Role, error)
	FindAll(ctx context.Context) ([]model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindRoleByName(ctx context.Context, roleName string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
