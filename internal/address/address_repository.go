package address

import (
	"context"
	"rijig/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	CreateAddress(ctx context.Context, address *model.Address) error
	FindAddressByUserID(ctx context.Context, userID string) ([]model.Address, error)
	FindAddressByID(ctx context.Context, id string) (*model.Address, error)
	UpdateAddress(ctx context.Context, address *model.Address) error
	DeleteAddress(ctx context.Context, id string) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db}
}

func (r *addressRepository) CreateAddress(ctx context.Context, address *model.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *addressRepository) FindAddressByUserID(ctx context.Context, userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) FindAddressByID(ctx context.Context, id string) (*model.Address, error) {
	var address model.Address
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, address *model.Address) error {
	err := r.db.WithContext(ctx).Save(address).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Address{}).Error
	if err != nil {
		return err
	}
	return nil
}
