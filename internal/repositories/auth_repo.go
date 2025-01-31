package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByIdentifierAndRole(identifier, roleID string) (*model.User, error)
	FindByEmailOrUsernameOrPhone(identifier string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByPhoneAndRole(phone, roleID string) (*model.User, error)
	FindByEmailAndRole(email, roleID string) (*model.User, error)

	Create(user *model.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByIdentifierAndRole(identifier, roleID string) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Role").Where("(email = ? OR username = ? OR phone = ?) AND role_id = ?", identifier, identifier, identifier, roleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
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

func (r *userRepository) FindByEmailAndRole(email, roleID string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ? AND role_id = ?", email, roleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
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
