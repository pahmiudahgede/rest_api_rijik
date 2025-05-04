package repositories

import (
	"rijig/model"

	"gorm.io/gorm"
)

type AuthPengepulRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetUserByPhoneAndRole(phone string, roleId string) (*model.User, error)
}

type authPengepulRepository struct {
	db *gorm.DB
}

func NewAuthPengepulRepositories(db *gorm.DB) AuthPengepulRepository {
	return &authPengepulRepository{db}
}

func (r *authPengepulRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authPengepulRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authPengepulRepository) GetUserByPhoneAndRole(phone string, roleId string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ? AND role_id = ?", phone, roleId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
