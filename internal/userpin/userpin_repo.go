package userpin

import (
	"context"
	"rijig/model"

	"gorm.io/gorm"
)

type UserPinRepository interface {
	FindByUserID(ctx context.Context, userID string) (*model.UserPin, error)
	Create(ctx context.Context, userPin *model.UserPin) error
	Update(ctx context.Context, userPin *model.UserPin) error
}

type userPinRepository struct {
	db *gorm.DB
}

func NewUserPinRepository(db *gorm.DB) UserPinRepository {
	return &userPinRepository{db}
}

func (r *userPinRepository) FindByUserID(ctx context.Context, userID string) (*model.UserPin, error) {
	var userPin model.UserPin
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userPin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &userPin, nil
}

func (r *userPinRepository) Create(ctx context.Context, userPin *model.UserPin) error {
	err := r.db.WithContext(ctx).Create(userPin).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userPinRepository) Update(ctx context.Context, userPin *model.UserPin) error {
	err := r.db.WithContext(ctx).Save(userPin).Error
	if err != nil {
		return err
	}
	return nil
}
