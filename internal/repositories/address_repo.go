package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type AddressRepository interface {
	CreateAddress(address *model.Address) error
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
