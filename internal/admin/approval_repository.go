package admin

import (
	"context"
	"fmt"
	"rijig/model"
	"strings"

	"gorm.io/gorm"
)

type ApprovalRepository interface {
	GetPendingUsers(ctx context.Context, req *GetPendingUsersRequest) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	UpdateUserRegistrationStatus(ctx context.Context, userID, status string, progress int8) error
	GetApprovalSummary(ctx context.Context) (*ApprovalSummary, error)
	GetUsersByIDs(ctx context.Context, userIDs []string) ([]model.User, error)
	BulkUpdateRegistrationStatus(ctx context.Context, userIDs []string, status string, progress int8) error
}

type approvalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{
		db: db,
	}
}

func (r *approvalRepository) GetPendingUsers(ctx context.Context, req *GetPendingUsersRequest) ([]model.User, int64, error) {
	var users []model.User
	var totalRecords int64

	query := r.db.WithContext(ctx).Model(&model.User{}).
		Preload("Role").
		Preload("IdentityCard").
		Preload("CompanyProfile")

	query = r.applyFilters(query, req)

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count pending users: %w", err)
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.
		Order("created_at DESC").
		Limit(req.Limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch pending users: %w", err)
	}

	return users, totalRecords, nil
}

func (r *approvalRepository) applyFilters(query *gorm.DB, req *GetPendingUsersRequest) *gorm.DB {

	if req.Status != "" {

		if req.Status == "pending" {
			req.Status = "awaiting_approval"
		}
		query = query.Where("registration_status = ?", req.Status)
	}

	if req.Role != "" {

		query = query.Joins("JOIN roles ON roles.id = users.role_id").
			Where("LOWER(roles.role_name) = ?", strings.ToLower(req.Role))
	}

	query = query.Where("registration_progress >= ?", 2)

	return query
}

func (r *approvalRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("IdentityCard").
		Preload("CompanyProfile").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *approvalRepository) UpdateUserRegistrationStatus(ctx context.Context, userID, status string, progress int8) error {
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"registration_status":   status,
			"registration_progress": progress,
			"updated_at":            "NOW()",
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update user registration status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *approvalRepository) GetApprovalSummary(ctx context.Context) (*ApprovalSummary, error) {
	var summary ApprovalSummary

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("registration_status = ? AND registration_progress >= ?", "awaiting_approval", 2).
		Count(&summary.TotalPending).Error; err != nil {
		return nil, fmt.Errorf("failed to count total pending users: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.registration_status = ? AND users.registration_progress >= ? AND LOWER(roles.role_name) = ?",
			"awaiting_approval", 2, "pengelola").
		Count(&summary.PengelolaPending).Error; err != nil {
		return nil, fmt.Errorf("failed to count pengelola pending: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.registration_status = ? AND users.registration_progress >= ? AND LOWER(roles.role_name) = ?",
			"awaiting_approval", 2, "pengepul").
		Count(&summary.PengepulPending).Error; err != nil {
		return nil, fmt.Errorf("failed to count pengepul pending: %w", err)
	}

	return &summary, nil
}

func (r *approvalRepository) GetUsersByIDs(ctx context.Context, userIDs []string) ([]model.User, error) {
	var users []model.User

	if err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("IdentityCard").
		Preload("CompanyProfile").
		Where("id IN ?", userIDs).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	return users, nil
}

func (r *approvalRepository) BulkUpdateRegistrationStatus(ctx context.Context, userIDs []string, status string, progress int8) error {
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id IN ?", userIDs).
		Updates(map[string]interface{}{
			"registration_status":   status,
			"registration_progress": progress,
			"updated_at":            "NOW()",
		})

	if result.Error != nil {
		return fmt.Errorf("failed to bulk update registration status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no users found to update")
	}

	return nil
}
