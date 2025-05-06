package repositories

import (
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type UserProfilRepository interface {
	FindByID(userID string) (*model.User, error)
	FindAll(page, limit int) ([]model.User, error)
	Update(user *model.User) error
	UpdateAvatar(userID, avatarURL string) error
	UpdatePassword(userID string, newPassword string) error
}

type userProfilRepository struct {
	DB *gorm.DB
}

func NewUserProfilRepository(db *gorm.DB) UserProfilRepository {
	return &userProfilRepository{DB: db}
}

func (r *userProfilRepository) FindByID(userID string) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Role").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return nil, fmt.Errorf("error finding user with ID %s: %v", userID, err)
	}

	if user.Role == nil {
		return nil, fmt.Errorf("role not found for user ID %s", userID)
	}

	return &user, nil
}

func (r *userProfilRepository) FindAll(page, limit int) ([]model.User, error) {
	var users []model.User
	offset := (page - 1) * limit
	err := r.DB.Preload("Role").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("error finding all users: %v", err)
	}
	return users, nil
}

func (r *userProfilRepository) Update(user *model.User) error {
	err := r.DB.Save(user).Error
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

func (r *userProfilRepository) UpdateAvatar(userID, avatarURL string) error {
	var user model.User
	err := r.DB.Model(&user).Where("id = ?", userID).Update("avatar", avatarURL).Error
	if err != nil {
		return fmt.Errorf("error updating avatar for user ID %s: %v", userID, err)
	}
	return nil
}

func (r *userProfilRepository) UpdatePassword(userID string, newPassword string) error {
	var user model.User
	err := r.DB.Model(&user).Where("id = ?", userID).Update("password", newPassword).Error
	if err != nil {
		return fmt.Errorf("error updating password for user ID %s: %v", userID, err)
	}
	return nil
}
