package repositories

import (
	"rijig/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id string) (*model.Role, error)
	FindRoleByName(roleName string) (*model.Role, error)
	FindAll() ([]model.Role, error)
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{DB: db}
}

func (r *roleRepository) FindByID(id string) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindRoleByName(roleName string) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
