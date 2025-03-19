package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByPhone(phone string) (*model.User, error)
	FindByPhoneAndRole(phone, roleID string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User

	err := r.DB.Preload("Role").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhoneAndRole(phone, roleID string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("phone = ? AND role_id = ?", phone, roleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}
