package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	FindByID(userID string) (*model.User, error)
	Update(user *model.User) error
}

type userProfileRepository struct {
	DB *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{DB: db}
}

func (r *userProfileRepository) FindByID(userID string) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Role").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userProfileRepository) Update(user *model.User) error {
	err := r.DB.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}
