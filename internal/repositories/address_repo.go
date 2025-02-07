package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type AddressRepository interface {
	CreateAddress(address *model.Address) error
	FindAddressByUserID(userID string) ([]model.Address, error)
	FindAddressByID(id string) (*model.Address, error)
	UpdateAddress(address *model.Address) error
}

type addressRepository struct {
	DB *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{DB: db}
}

func (r *addressRepository) CreateAddress(address *model.Address) error {
	return r.DB.Create(address).Error
}


func (r *addressRepository) FindAddressByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.DB.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) FindAddressByID(id string) (*model.Address, error) {
	var address model.Address
	err := r.DB.Where("id = ?", id).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) UpdateAddress(address *model.Address) error {
	err := r.DB.Save(address).Error
	if err != nil {
		return err
	}
	return nil
}