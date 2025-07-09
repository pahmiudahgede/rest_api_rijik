package admin

import (
	"context"
	"fmt"
	"math"
	"rijig/model"
	"rijig/utils"
	"strings"
	"time"
)

type ApprovalService interface {
	GetPendingUsers(ctx context.Context, req *GetPendingUsersRequest) (*PendingUsersListResponse, error)
	ApproveUser(ctx context.Context, userID, adminID string, notes string) (*ApprovalActionResponse, error)
	RejectUser(ctx context.Context, userID, adminID string, notes string) (*ApprovalActionResponse, error)
	ProcessApprovalAction(ctx context.Context, req *ApprovalActionRequest, adminID string) (*ApprovalActionResponse, error)
	BulkProcessApproval(ctx context.Context, req *BulkApprovalRequest, adminID string) (*BulkApprovalResponse, error)
	GetUserApprovalDetails(ctx context.Context, userID string) (*PendingUserResponse, error)
}

type approvalService struct {
	repo ApprovalRepository
}

func NewApprovalService(repo ApprovalRepository) ApprovalService {
	return &approvalService{
		repo: repo,
	}
}

func (s *approvalService) GetPendingUsers(ctx context.Context, req *GetPendingUsersRequest) (*PendingUsersListResponse, error) {

	req.SetDefaults()

	users, totalRecords, err := s.repo.GetPendingUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending users: %w", err)
	}

	summary, err := s.repo.GetApprovalSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get approval summary: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(req.Limit)))
	pagination := PaginationInfo{
		Page:         req.Page,
		Limit:        req.Limit,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      req.Page < totalPages,
		HasPrev:      req.Page > 1,
	}

	userResponses := make([]PendingUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.convertToUserResponse(user)
	}

	return &PendingUsersListResponse{
		Users:      userResponses,
		Pagination: pagination,
		Summary:    *summary,
	}, nil
}

func (s *approvalService) ApproveUser(ctx context.Context, userID, adminID string, notes string) (*ApprovalActionResponse, error) {

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.validateUserForApproval(user, "approved"); err != nil {
		return nil, err
	}

	previousStatus := user.RegistrationStatus

	newStatus := utils.RegStatusConfirmed
	newProgress := utils.ProgressDataSubmitted

	if err := s.repo.UpdateUserRegistrationStatus(ctx, userID, newStatus, int8(newProgress)); err != nil {
		return nil, fmt.Errorf("failed to approved user: %w", err)
	}

	if err := s.revokeUserTokens(userID); err != nil {

		fmt.Printf("Warning: failed to revoke tokens for user %s: %v\n", userID, err)
	}

	return &ApprovalActionResponse{
		UserID:         userID,
		Action:         "approved",
		PreviousStatus: previousStatus,
		NewStatus:      newStatus,
		ProcessedAt:    time.Now(),
		ProcessedBy:    adminID,
		Notes:          notes,
	}, nil
}

func (s *approvalService) RejectUser(ctx context.Context, userID, adminID string, notes string) (*ApprovalActionResponse, error) {

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.validateUserForApproval(user, "rejected"); err != nil {
		return nil, err
	}

	previousStatus := user.RegistrationStatus

	newStatus := utils.RegStatusRejected
	newProgress := utils.ProgressOTPVerified

	if err := s.repo.UpdateUserRegistrationStatus(ctx, userID, newStatus, int8(newProgress)); err != nil {
		return nil, fmt.Errorf("failed to rejected user: %w", err)
	}

	if err := s.revokeUserTokens(userID); err != nil {

		fmt.Printf("Warning: failed to revoke tokens for user %s: %v\n", userID, err)
	}

	return &ApprovalActionResponse{
		UserID:         userID,
		Action:         "rejected",
		PreviousStatus: previousStatus,
		NewStatus:      newStatus,
		ProcessedAt:    time.Now(),
		ProcessedBy:    adminID,
		Notes:          notes,
	}, nil
}

func (s *approvalService) ProcessApprovalAction(ctx context.Context, req *ApprovalActionRequest, adminID string) (*ApprovalActionResponse, error) {
	switch req.Action {
	case "approved":
		return s.ApproveUser(ctx, req.UserID, adminID, req.Notes)
	case "rejected":
		return s.RejectUser(ctx, req.UserID, adminID, req.Notes)
	default:
		return nil, fmt.Errorf("invalid action: %s", req.Action)
	}
}

func (s *approvalService) BulkProcessApproval(ctx context.Context, req *BulkApprovalRequest, adminID string) (*BulkApprovalResponse, error) {

	users, err := s.repo.GetUsersByIDs(ctx, req.UserIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var results []ApprovalActionResponse
	var failures []ApprovalFailure
	successCount := 0

	for _, user := range users {

		if err := s.validateUserForApproval(&user, req.Action); err != nil {
			failures = append(failures, ApprovalFailure{
				UserID: user.ID,
				Error:  "validation_failed",
				Reason: err.Error(),
			})
			continue
		}

		actionReq := &ApprovalActionRequest{
			UserID: user.ID,
			Action: req.Action,
			Notes:  req.Notes,
		}

		result, err := s.ProcessApprovalAction(ctx, actionReq, adminID)
		if err != nil {
			failures = append(failures, ApprovalFailure{
				UserID: user.ID,
				Error:  "processing_failed",
				Reason: err.Error(),
			})
			continue
		}

		results = append(results, *result)
		successCount++
	}

	return &BulkApprovalResponse{
		SuccessCount: successCount,
		FailureCount: len(failures),
		Results:      results,
		Failures:     failures,
	}, nil
}

func (s *approvalService) GetUserApprovalDetails(ctx context.Context, userID string) (*PendingUserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	userResponse := s.convertToUserResponse(*user)
	return &userResponse, nil
}

func (s *approvalService) validateUserForApproval(user *model.User, action string) error {

	if user.Role == nil {
		return fmt.Errorf("user role not found")
	}

	roleName := strings.ToLower(user.Role.RoleName)
	if roleName != "pengelola" && roleName != "pengepul" {
		return fmt.Errorf("only pengelola and pengepul can be approved/rejected")
	}

	if user.RegistrationStatus != utils.RegStatusPending {
		return fmt.Errorf("user is not in awaiting_approval status")
	}

	if user.RegistrationProgress < utils.ProgressDataSubmitted {
		return fmt.Errorf("user has not submitted required data yet")
	}

	if roleName == "pengepul" && user.IdentityCard == nil {
		return fmt.Errorf("pengepul must have identity card data")
	}

	if roleName == "pengelola" && user.CompanyProfile == nil {
		return fmt.Errorf("pengelola must have company profile data")
	}

	return nil
}

func (s *approvalService) revokeUserTokens(userID string) error {

	return utils.RevokeAllRefreshTokens(userID)
}

func (s *approvalService) convertToUserResponse(user model.User) PendingUserResponse {
	response := PendingUserResponse{
		ID:                   user.ID,
		Name:                 user.Name,
		Phone:                user.Phone,
		Email:                user.Email,
		RegistrationStatus:   user.RegistrationStatus,
		RegistrationProgress: user.RegistrationProgress,
		SubmittedAt:          user.UpdatedAt,
	}

	if user.Role != nil {
		response.Role = RoleInfo{
			ID:       user.Role.ID,
			RoleName: user.Role.RoleName,
		}
	}

	stepInfo := utils.GetRegistrationStepInfo(
		user.Role.RoleName,
		int(user.RegistrationProgress),
		user.RegistrationStatus,
	)
	if stepInfo != nil {
		response.RegistrationStepInfo = &RegistrationStepResponse{
			Step:                  stepInfo.Step,
			Status:                stepInfo.Status,
			Description:           stepInfo.Description,
			RequiresAdminApproval: stepInfo.RequiresAdminApproval,
			IsAccessible:          stepInfo.IsAccessible,
			IsCompleted:           stepInfo.IsCompleted,
		}
	}

	if user.IdentityCard != nil {
		response.IdentityCard = &IdentityCardInfo{
			ID:                   user.IdentityCard.ID,
			IdentificationNumber: user.IdentityCard.Identificationumber,
			Fullname:             user.IdentityCard.Fullname,
			Placeofbirth:         user.IdentityCard.Placeofbirth,
			Dateofbirth:          user.IdentityCard.Dateofbirth,
			Gender:               user.IdentityCard.Gender,
			BloodType:            user.IdentityCard.BloodType,
			Province:             user.IdentityCard.Province,
			District:             user.IdentityCard.District,
			SubDistrict:          user.IdentityCard.SubDistrict,
			Village:              user.IdentityCard.Village,
			PostalCode:           user.IdentityCard.PostalCode,
			Religion:             user.IdentityCard.Religion,
			Maritalstatus:        user.IdentityCard.Maritalstatus,
			Job:                  user.IdentityCard.Job,
			Citizenship:          user.IdentityCard.Citizenship,
			Validuntil:           user.IdentityCard.Validuntil,
			Cardphoto:            user.IdentityCard.Cardphoto,
		}
	}

	if user.CompanyProfile != nil {
		response.CompanyProfile = &CompanyProfileInfo{
			ID:                 user.CompanyProfile.ID,
			CompanyName:        user.CompanyProfile.CompanyName,
			CompanyAddress:     user.CompanyProfile.CompanyAddress,
			CompanyPhone:       user.CompanyProfile.CompanyPhone,
			CompanyEmail:       user.CompanyProfile.CompanyEmail,
			CompanyLogo:        user.CompanyProfile.CompanyLogo,
			CompanyWebsite:     user.CompanyProfile.CompanyWebsite,
			TaxID:              user.CompanyProfile.TaxID,
			FoundedDate:        user.CompanyProfile.FoundedDate,
			CompanyType:        user.CompanyProfile.CompanyType,
			CompanyDescription: user.CompanyProfile.CompanyDescription,
		}
	}

	return response
}
