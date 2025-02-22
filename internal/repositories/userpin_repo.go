package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type UserPinRepository interface {
	FindByUserID(userID string) (*model.UserPin, error)
	FindByPin(userPin string) (*model.UserPin, error)
	Create(userPin *model.UserPin) error
	Update(userPin *model.UserPin) error
}

type userPinRepository struct {
	DB *gorm.DB
}

func NewUserPinRepository(db *gorm.DB) UserPinRepository {
	return &userPinRepository{DB: db}
}

func (r *userPinRepository) FindByUserID(userID string) (*model.UserPin, error) {
	var userPin model.UserPin
	err := r.DB.Where("user_id = ?", userID).First(&userPin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			return nil, nil
		}
		return nil, err
	}
	return &userPin, nil
}

func (r *userPinRepository) FindByPin(pin string) (*model.UserPin, error) {
	var userPin model.UserPin
	err := r.DB.Where("pin = ?", pin).First(&userPin).Error
	if err != nil {
		return nil, err
	}
	return &userPin, nil
}

func (r *userPinRepository) Create(userPin *model.UserPin) error {
	err := r.DB.Create(userPin).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userPinRepository) Update(userPin *model.UserPin) error {
	err := r.DB.Save(userPin).Error
	if err != nil {
		return err
	}
	return nil
}
