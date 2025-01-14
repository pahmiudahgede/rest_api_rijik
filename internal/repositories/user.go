package repositories

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/config"
)

func GetUsers() ([]domain.User, error) {
	var users []domain.User

	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByRole(roleID string) ([]domain.User, error) {
	var users []domain.User

	if err := config.DB.Where("role_id = ?", roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByUserrId(userID string) (domain.User, error) {
	var user domain.User

	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}