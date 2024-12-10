package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetUserRoleByID(id string) (domain.UserRole, error) {
	var role domain.UserRole
	err := config.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return role, errors.New("userRole tidak ditemukan")
	}
	return role, nil
}

func GetAllUserRoles() ([]domain.UserRole, error) {
	var roles []domain.UserRole
	err := config.DB.Find(&roles).Error
	if err != nil {
		return nil, errors.New("gagal mengambil data UserRole")
	}
	return roles, nil
}
