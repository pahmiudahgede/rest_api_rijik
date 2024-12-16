package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetRequestPickupsByUser(userID string) ([]domain.RequestPickup, error) {
	var requestPickups []domain.RequestPickup

	err := config.DB.Preload("Request").
		Preload("Request.TrashCategory").
		Preload("UserAddress").
		Where("user_id = ?", userID).
		Find(&requestPickups).Error

	if err != nil {
		return nil, err
	}

	return requestPickups, nil
}
