package repositories

import (
	"fmt"

	"rijig/model"

	"gorm.io/gorm"
)

type StoreRepository interface {
	FindStoreByUserID(userID string) (*model.Store, error)
	FindStoreByID(storeID string) (*model.Store, error)
	FindAddressByID(addressID string) (*model.Address, error)

	CreateStore(store *model.Store) error
	UpdateStore(store *model.Store) error

	DeleteStore(storeID string) error
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

func (r *storeRepository) FindStoreByID(storeID string) (*model.Store, error) {
	var store model.Store
	if err := r.DB.Where("id = ?", storeID).First(&store).Error; err != nil {
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

func (r *storeRepository) UpdateStore(store *model.Store) error {
	if err := r.DB.Save(store).Error; err != nil {
		return err
	}
	return nil
}

func (r *storeRepository) DeleteStore(storeID string) error {

	if storeID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}

	if err := r.DB.Where("id = ?", storeID).Delete(&model.Store{}).Error; err != nil {
		return fmt.Errorf("failed to delete store: %w", err)
	}

	return nil
}
