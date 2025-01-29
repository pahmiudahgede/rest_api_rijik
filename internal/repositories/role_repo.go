package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id string) (*model.Role, error)
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
