package admin

import (
	"context"
	"fmt"
	"log"
)

type AdminService interface {
	GetAllUsers(ctx context.Context, req GetAllUsersRequest) (interface{}, *int, *int, int64, error)
	UpdateRegistrationStatus(ctx context.Context, userID string, req UpdateRegistrationStatusRequest) error
	GetUserStatistics(ctx context.Context) (map[string]interface{}, error)
	ValidatePaginationParams(page, limit *int) (*int, *int, error)
	GetMessage(role string, hasPagination bool, total int64) string
}

type adminService struct {
	adminRepo AdminRepository
}

func NewAdminService(adminRepo AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (s *adminService) GetAllUsers(ctx context.Context, req GetAllUsersRequest) (interface{}, *int, *int, int64, error) {

	if err := s.validateRole(req.Role); err != nil {
		return nil, nil, nil, 0, err
	}

	result, err := s.adminRepo.GetAllUsers(ctx, req)
	if err != nil {
		log.Printf("Error fetching users from repository: %v", err)
		return nil, nil, nil, 0, fmt.Errorf("failed to fetch users")
	}

	responseData, err := s.buildRoleBasedResponse(result.Users, req.Role)
	if err != nil {
		log.Printf("Error building role-based response: %v", err)
		return nil, nil, nil, 0, fmt.Errorf("failed to build response")
	}

	return responseData, req.Page, req.Limit, result.Total, nil
}

func (s *adminService) UpdateRegistrationStatus(ctx context.Context, userID string, req UpdateRegistrationStatusRequest) error {

	user, err := s.adminRepo.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Error finding user %s: %v", userID, err)
		return fmt.Errorf("user not found")
	}

	if err := s.validateRegistrationStatusUpdate(user.RegistrationStatus, req.Action); err != nil {
		return err
	}

	if err := s.adminRepo.UpdateRegistrationStatus(ctx, userID, req.Action); err != nil {
		log.Printf("Error updating registration status for user %s: %v", userID, err)
		return fmt.Errorf("failed to update registration status")
	}

	log.Printf("Successfully updated registration status for user %s to %s", userID, req.Action)
	return nil
}

func (s *adminService) validateRole(role string) error {
	validRoles := map[string]bool{
		"masyarakat": true,
		"pengepul":   true,
		"pengelola":  true,
	}

	if !validRoles[role] {
		return fmt.Errorf("invalid role: %s", role)
	}

	return nil
}

func (s *adminService) validateRegistrationStatusUpdate(currentStatus, action string) error {

	if currentStatus != "awaiting_approval" && currentStatus != "pending" {
		return fmt.Errorf("cannot update registration status: user is already %s", currentStatus)
	}

	validActions := map[string]bool{
		"approved": true,
		"rejected": true,
	}

	if !validActions[action] {
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
}

func (s *adminService) buildRoleBasedResponse(userRelations []UserWithRelations, role string) (interface{}, error) {
	switch role {
	case "masyarakat":
		return s.buildMasyarakatResponse(userRelations), nil
	case "pengepul":
		return s.buildPengepulResponse(userRelations), nil
	case "pengelola":
		return s.buildPengelolaResponse(userRelations), nil
	default:
		return nil, fmt.Errorf("unsupported role: %s", role)
	}
}

func (s *adminService) buildMasyarakatResponse(userRelations []UserWithRelations) []MasyarakatUserResponse {
	responses := make([]MasyarakatUserResponse, 0, len(userRelations))

	for _, userRelation := range userRelations {
		response := ToMasyarakatResponse(userRelation)
		responses = append(responses, response)
	}

	return responses
}

func (s *adminService) buildPengepulResponse(userRelations []UserWithRelations) []PengepulUserResponse {
	responses := make([]PengepulUserResponse, 0, len(userRelations))

	for _, userRelation := range userRelations {
		response := ToPengepulResponse(userRelation)
		responses = append(responses, response)
	}

	return responses
}

func (s *adminService) buildPengelolaResponse(userRelations []UserWithRelations) []PengelolaUserResponse {
	responses := make([]PengelolaUserResponse, 0, len(userRelations))

	for _, userRelation := range userRelations {
		response := ToPengelolaResponse(userRelation)
		responses = append(responses, response)
	}

	return responses
}

func (s *adminService) GetUserStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	roles := []string{"masyarakat", "pengepul", "pengelola"}
	for _, role := range roles {
		req := GetAllUsersRequest{Role: role}
		result, err := s.adminRepo.GetAllUsers(ctx, req)
		if err != nil {
			log.Printf("Error getting statistics for role %s: %v", role, err)
			continue
		}
		stats[fmt.Sprintf("total_%s", role)] = result.Total
	}

	pendingReq := GetAllUsersRequest{
		Role:      "pengepul",
		StatusReg: "awaiting_approval",
	}
	pendingPengepul, err := s.adminRepo.GetAllUsers(ctx, pendingReq)
	if err == nil {
		stats["pending_pengepul_approvals"] = pendingPengepul.Total
	}

	pendingReq.Role = "pengelola"
	pendingPengelola, err := s.adminRepo.GetAllUsers(ctx, pendingReq)
	if err == nil {
		stats["pending_pengelola_approvals"] = pendingPengelola.Total
	}

	return stats, nil
}

func (s *adminService) ValidatePaginationParams(page, limit *int) (*int, *int, error) {

	if page == nil && limit == nil {
		return nil, nil, nil
	}

	if page == nil || limit == nil {
		return nil, nil, fmt.Errorf("both page and limit must be provided for pagination")
	}

	if *page < 1 {
		return nil, nil, fmt.Errorf("page must be greater than 0")
	}

	if *limit < 1 {
		return nil, nil, fmt.Errorf("limit must be greater than 0")
	}

	maxLimit := 100
	if *limit > maxLimit {
		return nil, nil, fmt.Errorf("limit cannot exceed %d", maxLimit)
	}

	return page, limit, nil
}

func (s *adminService) GetMessage(role string, hasPagination bool, total int64) string {
	if hasPagination {
		return fmt.Sprintf("Successfully retrieved %s users with pagination", role)
	}

	if total == 0 {
		return fmt.Sprintf("No %s users found", role)
	}

	return fmt.Sprintf("Successfully retrieved all %s users", role)
}
