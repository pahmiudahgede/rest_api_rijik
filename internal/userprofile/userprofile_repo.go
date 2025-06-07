package userprofile

import (
	"context"
	"rijig/model"

	"gorm.io/gorm"
)

type AuthenticationRepository interface {
	UpdateUser(ctx context.Context, user *model.User) error
	PatchUser(ctx context.Context, userID string, updates map[string]interface{}) error
}

type authenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{db}
}

func (r *authenticationRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(user).Error
}

func (r *authenticationRepository) PatchUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}
