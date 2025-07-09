package admin

import "time"

// Request DTOs
type GetPendingUsersRequest struct {
	Role   string `query:"role" validate:"omitempty,oneof=pengelola pengepul"`
	Status string `query:"status" validate:"omitempty,oneof=awaiting_approval pending"`
	Page   int    `query:"page" validate:"min=1"`
	Limit  int    `query:"limit" validate:"min=1,max=100"`
}

type ApprovalActionRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Notes  string `json:"notes" validate:"omitempty,max=500"`
}

type BulkApprovalRequest struct {
	UserIDs []string `json:"user_ids" validate:"required,min=1,max=50,dive,uuid"`
	Action  string   `json:"action" validate:"required,oneof=approve reject"`
	Notes   string   `json:"notes" validate:"omitempty,max=500"`
}

// Response DTOs
type PendingUserResponse struct {
	ID                   string                    `json:"id"`
	Name                 string                    `json:"name,omitempty"`
	Phone                string                    `json:"phone"`
	Email                string                    `json:"email,omitempty"`
	Role                 RoleInfo                  `json:"role"`
	RegistrationStatus   string                    `json:"registration_status"`
	RegistrationProgress int8                      `json:"registration_progress"`
	SubmittedAt          time.Time                 `json:"submitted_at"`
	IdentityCard         *IdentityCardInfo         `json:"identity_card,omitempty"`
	CompanyProfile       *CompanyProfileInfo       `json:"company_profile,omitempty"`
	RegistrationStepInfo *RegistrationStepResponse `json:"step_info"`
}

type RoleInfo struct {
	ID       string `json:"id"`
	RoleName string `json:"role_name"`
}

type IdentityCardInfo struct {
	ID                   string `json:"id"`
	IdentificationNumber string `json:"identification_number"`
	Fullname             string `json:"fullname"`
	Placeofbirth         string `json:"place_of_birth"`
	Dateofbirth          string `json:"date_of_birth"`
	Gender               string `json:"gender"`
	BloodType            string `json:"blood_type"`
	Province             string `json:"province"`
	District             string `json:"district"`
	SubDistrict          string `json:"sub_district"`
	Village              string `json:"village"`
	PostalCode           string `json:"postal_code"`
	Religion             string `json:"religion"`
	Maritalstatus        string `json:"marital_status"`
	Job                  string `json:"job"`
	Citizenship          string `json:"citizenship"`
	Validuntil           string `json:"valid_until"`
	Cardphoto            string `json:"card_photo"`
}

type CompanyProfileInfo struct {
	ID                 string `json:"id"`
	CompanyName        string `json:"company_name"`
	CompanyAddress     string `json:"company_address"`
	CompanyPhone       string `json:"company_phone"`
	CompanyEmail       string `json:"company_email"`
	CompanyLogo        string `json:"company_logo"`
	CompanyWebsite     string `json:"company_website"`
	TaxID              string `json:"tax_id"`
	FoundedDate        string `json:"founded_date"`
	CompanyType        string `json:"company_type"`
	CompanyDescription string `json:"company_description"`
}

type RegistrationStepResponse struct {
	Step                  int    `json:"step"`
	Status                string `json:"status"`
	Description           string `json:"description"`
	RequiresAdminApproval bool   `json:"requires_admin_approval"`
	IsAccessible          bool   `json:"is_accessible"`
	IsCompleted           bool   `json:"is_completed"`
}

type PendingUsersListResponse struct {
	Users      []PendingUserResponse `json:"users"`
	Pagination PaginationInfo        `json:"pagination"`
	Summary    ApprovalSummary       `json:"summary"`
}

type PaginationInfo struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

type ApprovalSummary struct {
	TotalPending     int64 `json:"total_pending"`
	PengelolaPending int64 `json:"pengelola_pending"`
	PengepulPending  int64 `json:"pengepul_pending"`
}

type ApprovalActionResponse struct {
	UserID         string    `json:"user_id"`
	Action         string    `json:"action"`
	PreviousStatus string    `json:"previous_status"`
	NewStatus      string    `json:"new_status"`
	ProcessedAt    time.Time `json:"processed_at"`
	ProcessedBy    string    `json:"processed_by"`
	Notes          string    `json:"notes,omitempty"`
}

type BulkApprovalResponse struct {
	SuccessCount int                      `json:"success_count"`
	FailureCount int                      `json:"failure_count"`
	Results      []ApprovalActionResponse `json:"results"`
	Failures     []ApprovalFailure        `json:"failures,omitempty"`
}

type ApprovalFailure struct {
	UserID string `json:"user_id"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// Validation helper
func (r *GetPendingUsersRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 20
	}
	if r.Status == "" {
		r.Status = "awaiting_approval"
	}
}
