package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type StoreRepository interface {
	FindStoreByUserID(userID string) (*model.Store, error)
	FindAddressByID(addressID string) (*model.Address, error)
	CreateStore(store *model.Store) error
}

type storeRepository struct {
	DB *gorm.DB
}

func NewStoreRepository(DB *gorm.DB) StoreRepository {
	return &storeRepository{DB}
}

func (r *storeRepository) FindStoreByUserID(userID string) (*model.Store, error) {
	var store model.Store
	if err := r.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindAddressByID(addressID string) (*model.Address, error) {
	var address model.Address
	if err := r.DB.Where("id = ?", addressID).First(&address).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

func (r *storeRepository) CreateStore(store *model.Store) error {
	if err := r.DB.Create(store).Error; err != nil {
		return err
	}
	return nil
}
