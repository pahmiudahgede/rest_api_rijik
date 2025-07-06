package authentication

import (
	"context"
	"fmt"
	"log"
	"rijig/model"

	"gorm.io/gorm"
)

type AuthenticationRepository interface {
	FindUserByPhone(ctx context.Context, phone string) (*model.User, error)
	FindUserByPhoneAndRole(ctx context.Context, phone, rolename string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, userID string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	PatchUser(ctx context.Context, userID string, updates map[string]interface{}) error

	GetIdentityCardsByUserRegStatus(ctx context.Context, userRegStatus string) ([]model.IdentityCard, error)
	GetCompanyProfilesByUserRegStatus(ctx context.Context, userRegStatus string) ([]model.CompanyProfile, error)
}

type authenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{db}
}

func (r *authenticationRepository) FindUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authenticationRepository) FindUserByPhoneAndRole(ctx context.Context, phone, rolename string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Joins("JOIN roles AS role ON role.id = users.role_id").
		Where("users.phone = ? AND role.role_name = ?", phone, rolename).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authenticationRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authenticationRepository) FindUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authenticationRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
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

func (r *authenticationRepository) GetIdentityCardsByUserRegStatus(ctx context.Context, userRegStatus string) ([]model.IdentityCard, error) {
	var identityCards []model.IdentityCard

	if err := r.db.WithContext(ctx).
		Joins("JOIN users ON identity_cards.user_id = users.id").
		Where("users.registration_status = ?", userRegStatus).
		Preload("User").
		Preload("User.Role").
		Find(&identityCards).Error; err != nil {
		log.Printf("Error fetching identity cards by user registration status: %v", err)
		return nil, fmt.Errorf("error fetching identity cards by user registration status: %w", err)
	}

	log.Printf("Found %d identity cards with registration status: %s", len(identityCards), userRegStatus)
	return identityCards, nil
}

func (r *authenticationRepository) GetCompanyProfilesByUserRegStatus(ctx context.Context, userRegStatus string) ([]model.CompanyProfile, error) {
	var companyProfiles []model.CompanyProfile

	if err := r.db.WithContext(ctx).
		Joins("JOIN users ON company_profiles.user_id = users.id").
		Where("users.registration_status = ?", userRegStatus).
		Preload("User").
		Preload("User.Role").
		Find(&companyProfiles).Error; err != nil {
		log.Printf("Error fetching company profiles by user registration status: %v", err)
		return nil, fmt.Errorf("error fetching company profiles by user registration status: %w", err)
	}

	log.Printf("Found %d company profiles with registration status: %s", len(companyProfiles), userRegStatus)
	return companyProfiles, nil
}
