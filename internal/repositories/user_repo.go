package repositories

import (
	"fmt"

	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	FindByID(userID string) (*model.User, error)
	Update(user *model.User) error
	UpdateAvatar(userID, avatarURL string) error
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
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return nil, err
	}

	if user.Role == nil {
		return nil, fmt.Errorf("role not found for this user")
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

func (r *userProfileRepository) UpdateAvatar(userID, avatarURL string) error {
	var user model.User
	err := r.DB.Model(&user).Where("id = ?", userID).Update("avatar", avatarURL).Error
	if err != nil {
		return err
	}
	return nil
}
