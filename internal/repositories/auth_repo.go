package repositories

import (
	"rijig/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetUserByPhoneAndRole(phone string, roleID string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByPhoneAndRole(phone string, roleID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ? AND role_id = ?", phone, roleID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
