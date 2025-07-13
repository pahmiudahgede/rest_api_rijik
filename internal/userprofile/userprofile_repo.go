package userprofile

import (
	"context"
	"errors"
	"rijig/model"

	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetByID(ctx context.Context, userID string) (*model.User, error)
	GetByRoleName(ctx context.Context, roleName string) ([]*model.User, error)
	Update(ctx context.Context, userID string, user *model.User) error
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

func (r *userProfileRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userProfileRepository) GetByRoleName(ctx context.Context, roleName string) ([]*model.User, error) {
	var users []*model.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.role_name = ?", roleName).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userProfileRepository) Update(ctx context.Context, userID string, user *model.User) error {
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
