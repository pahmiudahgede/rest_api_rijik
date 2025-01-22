package repositories

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetUserRoleByID(id string) (domain.UserRole, error) {
	var role domain.UserRole

	ctx := config.Context()
	cachedRole, err := config.RedisClient.Get(ctx, "userRole:"+id).Result()

	if err == nil {

		err := json.Unmarshal([]byte(cachedRole), &role)
		if err != nil {
			return role, errors.New("gagal mendekode data cache Redis")
		}
		return role, nil
	}

	err = config.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return role, errors.New("userRole tidak ditemukan")
	}

	roleJSON, err := json.Marshal(role)
	if err != nil {
		return role, errors.New("gagal mendekode data untuk Redis")
	}

	err = config.RedisClient.Set(ctx, "userRole:"+id, roleJSON, time.Hour*24).Err()
	if err != nil {
		return role, errors.New("gagal menyimpan data di Redis")
	}

	return role, nil
}

func GetAllUserRoles() ([]domain.UserRole, error) {
	var roles []domain.UserRole

	ctx := config.Context()
	cachedRoles, err := config.RedisClient.Get(ctx, "allUserRoles").Result()

	if err == nil {

		err := json.Unmarshal([]byte(cachedRoles), &roles)
		if err != nil {
			return roles, errors.New("gagal mendekode data cache Redis")
		}
		return roles, nil
	}

	err = config.DB.Find(&roles).Error
	if err != nil {
		return nil, errors.New("gagal mengambil data UserRole")
	}

	rolesJSON, err := json.Marshal(roles)
	if err != nil {
		return roles, errors.New("gagal mendekode data untuk Redis")
	}

	err = config.RedisClient.Set(ctx, "allUserRoles", rolesJSON, time.Hour*24).Err()
	if err != nil {
		return roles, errors.New("gagal menyimpan data di Redis")
	}

	return roles, nil
}
