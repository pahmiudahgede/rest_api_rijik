package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetUserRoleByID(id string) (domain.UserRole, error) {
	role, err := repositories.GetUserRoleByID(id)
	if err != nil {
		return role, errors.New("UserRole tidak ditemukan")
	}
	return role, nil
}

func GetAllUserRoles() ([]domain.UserRole, error) {
	roles, err := repositories.GetAllUserRoles()
	if err != nil {
		return nil, errors.New("Gagal mengambil data UserRole")
	}
	return roles, nil
}
