package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetStoreByID(storeID string) (domain.Store, error) {
	var store domain.Store
	err := config.DB.Where("id = ?", storeID).First(&store).Error
	return store, err
}

func GetStoresByUserID(userID string, limit, offset int) ([]domain.Store, error) {
	var stores []domain.Store
	query := config.DB.Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&stores).Error
	return stores, err
}
