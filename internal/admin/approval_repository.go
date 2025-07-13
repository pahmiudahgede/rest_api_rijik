package admin

import (
	"context"
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetAllUsers(ctx context.Context, req GetAllUsersRequest) (*PaginatedUsersResult, error)
	UpdateRegistrationStatus(ctx context.Context, userID string, status string) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{
		db: db,
	}
}

func (r *adminRepository) GetAllUsers(ctx context.Context, req GetAllUsersRequest) (*PaginatedUsersResult, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("roles.role_name = ?", req.Role)

	if req.StatusReg != "" {
		query = query.Where("users.registration_status = ?", req.StatusReg)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	if req.Page != nil && req.Limit != nil {
		offset := (*req.Page - 1) * *req.Limit
		query = query.Offset(offset).Limit(*req.Limit)
	}

	query = query.Preload("Role")

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	userRelations, err := r.fetchUserRelations(ctx, users, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user relations: %w", err)
	}

	return &PaginatedUsersResult{
		Users: userRelations,
		Total: total,
	}, nil
}

func (r *adminRepository) UpdateRegistrationStatus(ctx context.Context, userID string, status string) error {

	var registrationStatus string
	switch status {
	case "approved":
		registrationStatus = "approved"
	case "rejected":
		registrationStatus = "rejected"
	default:
		return fmt.Errorf("invalid action: %s", status)
	}

	result := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("registration_status", registrationStatus)

	if result.Error != nil {
		return fmt.Errorf("failed to update registration status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", userID)
	}

	return nil
}

func (r *adminRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (r *adminRepository) fetchUserRelations(ctx context.Context, users []model.User, role string) ([]UserWithRelations, error) {
	userRelations := make([]UserWithRelations, 0, len(users))

	switch role {
	case "pengepul":

		identityCards, err := r.getIdentityCardsByUserIDs(ctx, r.extractUserIDs(users))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch identity cards: %w", err)
		}

		identityCardMap := make(map[string]*model.IdentityCard)
		for i := range identityCards {
			identityCardMap[identityCards[i].UserID] = &identityCards[i]
		}

		for _, user := range users {
			userRelations = append(userRelations, UserWithRelations{
				User:           user,
				IdentityCard:   identityCardMap[user.ID],
				CompanyProfile: nil,
			})
		}

	case "pengelola":

		companyProfiles, err := r.getCompanyProfilesByUserIDs(ctx, r.extractUserIDs(users))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch company profiles: %w", err)
		}

		companyProfileMap := make(map[string]*model.CompanyProfile)
		for i := range companyProfiles {
			companyProfileMap[companyProfiles[i].UserID] = &companyProfiles[i]
		}

		for _, user := range users {
			userRelations = append(userRelations, UserWithRelations{
				User:           user,
				IdentityCard:   nil,
				CompanyProfile: companyProfileMap[user.ID],
			})
		}

	case "masyarakat":

		for _, user := range users {
			userRelations = append(userRelations, UserWithRelations{
				User:           user,
				IdentityCard:   nil,
				CompanyProfile: nil,
			})
		}

	default:
		return nil, fmt.Errorf("unsupported role: %s", role)
	}

	return userRelations, nil
}

func (r *adminRepository) extractUserIDs(users []model.User) []string {
	userIDs := make([]string, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	return userIDs
}

func (r *adminRepository) getIdentityCardsByUserIDs(ctx context.Context, userIDs []string) ([]model.IdentityCard, error) {
	var identityCards []model.IdentityCard

	if len(userIDs) == 0 {
		return identityCards, nil
	}

	err := r.db.WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Find(&identityCards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch identity cards: %w", err)
	}

	return identityCards, nil
}

func (r *adminRepository) getCompanyProfilesByUserIDs(ctx context.Context, userIDs []string) ([]model.CompanyProfile, error) {
	var companyProfiles []model.CompanyProfile

	if len(userIDs) == 0 {
		return companyProfiles, nil
	}

	err := r.db.WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Find(&companyProfiles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch company profiles: %w", err)
	}

	return companyProfiles, nil
}

func (r *adminRepository) GetUsersByRole(ctx context.Context, role string) ([]model.User, error) {
	var users []model.User

	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("roles.role_name = ?", role).
		Preload("Role").
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users by role: %w", err)
	}

	return users, nil
}

func (r *adminRepository) GetUsersByStatus(ctx context.Context, status string) ([]model.User, error) {
	var users []model.User

	err := r.db.WithContext(ctx).
		Where("registration_status = ?", status).
		Preload("Role").
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users by status: %w", err)
	}

	return users, nil
}

func (r *adminRepository) BatchUpdateRegistrationStatus(ctx context.Context, userIDs []string, status string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id IN ?", userIDs).
		Update("registration_status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to batch update registration status: %w", result.Error)
	}

	return nil
}

func (r *adminRepository) GetIdentityCardByUserID(ctx context.Context, userID string) (*model.IdentityCard, error) {
	var identityCard model.IdentityCard

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&identityCard).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch identity card: %w", err)
	}

	return &identityCard, nil
}

func (r *adminRepository) GetCompanyProfileByUserID(ctx context.Context, userID string) (*model.CompanyProfile, error) {
	var companyProfile model.CompanyProfile

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&companyProfile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch company profile: %w", err)
	}

	return &companyProfile, nil
}
