package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmailOrUsernameOrPhone(identifier string) (*model.User, error)
	Create(user *model.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByEmailOrUsernameOrPhone(identifier string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ? OR username = ? OR phone = ?", identifier, identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	err := r.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
