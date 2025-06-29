package authentication

import (
	"rijig/utils"
	"strings"
	"time"
)

type LoginorRegistRequest struct {
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	RoleName string `json:"role_name"`
}

type VerifyOtpRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
	RoleName string `json:"role_name" validate:"required,oneof=masyarakat pengepul pengelola"`
	Phone    string `json:"phone" validate:"required"`
	Otp      string `json:"otp" validate:"required,len=6"`
}

type CreatePINRequest struct {
	PIN        string `json:"pin" validate:"required,len=6,numeric"`
	ConfirmPIN string `json:"confirm_pin" validate:"required,len=6,numeric"`
}

type VerifyPINRequest struct {
	PIN      string `json:"pin" validate:"required,len=6,numeric"`
	DeviceID string `json:"device_id" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
	UserID       string `json:"user_id" validate:"required"`
}

type LogoutRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
}

type OTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
	Phone     string `json:"phone"`
}

type AuthResponse struct {
	Message            string        `json:"message"`
	AccessToken        string        `json:"access_token,omitempty"`
	RefreshToken       string        `json:"refresh_token,omitempty"`
	TokenType          string        `json:"token_type,omitempty"`
	ExpiresIn          int64         `json:"expires_in,omitempty"`
	User               *UserResponse `json:"user,omitempty"`
	RegistrationStatus string        `json:"registration_status,omitempty"`
	NextStep           string        `json:"next_step,omitempty"`
	SessionID          string        `json:"session_id,omitempty"`
}

type UserResponse struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Phone                string  `json:"phone"`
	Email                string  `json:"email,omitempty"`
	Role                 string  `json:"role"`
	RegistrationStatus   string  `json:"registration_status"`
	RegistrationProgress int8    `json:"registration_progress"`
	PhoneVerified        bool    `json:"phone_verified"`
	Avatar               *string `json:"avatar,omitempty"`
	Gender               string  `json:"gender,omitempty"`
	DateOfBirth          string  `json:"date_of_birth,omitempty"`
	PlaceOfBirth         string  `json:"place_of_birth,omitempty"`
}

type RegistrationStatusResponse struct {
	CurrentStep        int                `json:"current_step"`
	TotalSteps         int                `json:"total_steps"`
	CompletedSteps     []RegistrationStep `json:"completed_steps"`
	NextStep           *RegistrationStep  `json:"next_step,omitempty"`
	RegistrationStatus string             `json:"registration_status"`
	IsComplete         bool               `json:"is_complete"`
	RequiresApproval   bool               `json:"requires_approval"`
	ApprovalMessage    string             `json:"approval_message,omitempty"`
}

type RegistrationStep struct {
	StepNumber  int    `json:"step_number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsRequired  bool   `json:"is_required"`
	IsCompleted bool   `json:"is_completed"`
	IsActive    bool   `json:"is_active"`
}

type OTPData struct {
	Phone     string    `json:"phone"`
	OTP       string    `json:"otp"`
	UserID    string    `json:"user_id,omitempty"`
	Role      string    `json:"role"`
	RoleID    string    `json:"role_id,omitempty"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Attempts  int       `json:"attempts"`
}

type SessionData struct {
	UserID    string    `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
	IsActive  bool      `json:"is_active"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

type ApproveRegistrationRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Message string `json:"message,omitempty"`
}

type RejectRegistrationRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}

type PendingRegistrationResponse struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Phone             string           `json:"phone"`
	Role              string           `json:"role"`
	RegistrationData  RegistrationData `json:"registration_data"`
	SubmittedAt       time.Time        `json:"submitted_at"`
	DocumentsUploaded []DocumentInfo   `json:"documents_uploaded"`
}

type RegistrationData struct {
	KTPNumber       string `json:"ktp_number,omitempty"`
	KTPImage        string `json:"ktp_image,omitempty"`
	FullName        string `json:"full_name,omitempty"`
	Address         string `json:"address,omitempty"`
	BusinessName    string `json:"business_name,omitempty"`
	BusinessType    string `json:"business_type,omitempty"`
	BusinessAddress string `json:"business_address,omitempty"`
	BusinessPhone   string `json:"business_phone,omitempty"`
	TaxNumber       string `json:"tax_number,omitempty"`
	BusinessLicense string `json:"business_license,omitempty"`
}

type DocumentInfo struct {
	Type        string    `json:"type"`
	FileName    string    `json:"file_name"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Status      string    `json:"status"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
}

type AuthStatsResponse struct {
	TotalUsers           int64                 `json:"total_users"`
	ActiveUsers          int64                 `json:"active_users"`
	PendingRegistrations int64                 `json:"pending_registrations"`
	UsersByRole          map[string]int64      `json:"users_by_role"`
	RegistrationStats    RegistrationStatsData `json:"registration_stats"`
	LoginStats           LoginStatsData        `json:"login_stats"`
}

type RegistrationStatsData struct {
	TotalRegistrations    int64 `json:"total_registrations"`
	CompletedToday        int64 `json:"completed_today"`
	CompletedThisWeek     int64 `json:"completed_this_week"`
	CompletedThisMonth    int64 `json:"completed_this_month"`
	PendingApproval       int64 `json:"pending_approval"`
	RejectedRegistrations int64 `json:"rejected_registrations"`
}

type LoginStatsData struct {
	TotalLogins      int64 `json:"total_logins"`
	LoginsToday      int64 `json:"logins_today"`
	LoginsThisWeek   int64 `json:"logins_this_week"`
	LoginsThisMonth  int64 `json:"logins_this_month"`
	UniqueUsersToday int64 `json:"unique_users_today"`
	UniqueUsersWeek  int64 `json:"unique_users_week"`
	UniqueUsersMonth int64 `json:"unique_users_month"`
}

type PaginationRequest struct {
	Page   int    `json:"page" query:"page" validate:"min=1"`
	Limit  int    `json:"limit" query:"limit" validate:"min=1,max=100"`
	Sort   string `json:"sort" query:"sort"`
	Order  string `json:"order" query:"order" validate:"oneof=asc desc"`
	Search string `json:"search" query:"search"`
	Filter string `json:"filter" query:"filter"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type PaginatedResponse struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type SMSWebhookRequest struct {
	MessageID string `json:"message_id"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type RateLimitInfo struct {
	Limit      int           `json:"limit"`
	Remaining  int           `json:"remaining"`
	ResetTime  time.Time     `json:"reset_time"`
	RetryAfter time.Duration `json:"retry_after,omitempty"`
}

type StepResponse struct {
	UserID               string `json:"user_id"`
	Role                 string `json:"role"`
	RegistrationStatus   string `json:"registration_status"`
	RegistrationProgress int    `json:"registration_progress"`
	NextStep             string `json:"next_step"`
}

type RegisterAdminRequest struct {
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	DateOfBirth     string `json:"dateofbirth"`
	PlaceOfBirth    string `json:"placeofbirth"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type LoginAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

type VerifyAdminOTPRequest struct {
	Email    string `json:"email" validate:"required,email"`
	OTP      string `json:"otp" validate:"required,len=6,numeric"`
	DeviceID string `json:"device_id" validate:"required"`
}

type ResendAdminOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type OTPAdminResponse struct {
	Message       string        `json:"message"`
	Email         string        `json:"email"`
	ExpiresIn     time.Duration `json:"expires_in_seconds"`
	RemainingTime string        `json:"remaining_time"`
	CanResend     bool          `json:"can_resend"`
	MaxAttempts   int           `json:"max_attempts"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type ResetPasswordResponse struct {
	Message       string        `json:"message"`
	Email         string        `json:"email"`
	ExpiresIn     time.Duration `json:"expires_in_seconds"`
	RemainingTime string        `json:"remaining_time"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailVerificationResponse struct {
	Message       string        `json:"message"`
	Email         string        `json:"email"`
	ExpiresIn     time.Duration `json:"expires_in_seconds"`
	RemainingTime string        `json:"remaining_time"`
}


func (r *LoginorRegistRequest) ValidateLoginorRegistRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if !utils.IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "nomor harus dimulai 62.. dan 8-14 digit")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

func (r *VerifyOtpRequest) ValidateVerifyOtpRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)
	if len(strings.TrimSpace(r.DeviceID)) < 10 {
		errors["device_id"] = append(errors["device_id"], "Device ID must be at least 10 characters")
	}

	validRoles := map[string]bool{"masyarakat": true, "pengepul": true, "pengelola": true}
	if _, ok := validRoles[r.RoleName]; !ok {
		errors["role"] = append(errors["role"], "Role tidak valid, hanya masyarakat, pengepul, atau pengelola")
	}

	if !utils.IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "nomor harus dimulai 62.. dan 8-14 digit")
	}

	if len(r.Otp) != 4 || !utils.IsNumeric(r.Otp) {
		errors["otp"] = append(errors["otp"], "OTP must be 4-digit number")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

func (r *LoginAdminRequest) ValidateLoginAdminRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if !utils.IsValidEmail(r.Email) {
		errors["email"] = append(errors["email"], "Invalid email format")
	}

	if strings.TrimSpace(r.Password) == "" {
		errors["password"] = append(errors["password"], "Password is required")
	}

	if len(strings.TrimSpace(r.DeviceID)) < 10 {
		errors["device_id"] = append(errors["device_id"], "Device ID must be at least 10 characters")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

func (r *RegisterAdminRequest) ValidateRegisterAdminRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	}

	genderLower := strings.ToLower(strings.TrimSpace(r.Gender))
	if genderLower != "laki-laki" && genderLower != "perempuan" {
		errors["gender"] = append(errors["gender"], "Gender must be either 'laki-laki' or 'perempuan'")
	}

	if strings.TrimSpace(r.DateOfBirth) == "" {
		errors["dateofbirth"] = append(errors["dateofbirth"], "Date of birth is required")
	} else {
		_, err := time.Parse("02-01-2006", r.DateOfBirth)
		if err != nil {
			errors["dateofbirth"] = append(errors["dateofbirth"], "Date of birth must be in DD-MM-YYYY format")
		}
	}

	if strings.TrimSpace(r.PlaceOfBirth) == "" {
		errors["placeofbirth"] = append(errors["placeofbirth"], "Place of birth is required")
	}

	if !utils.IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Phone must be valid, has 8-14 digit and start with '62..'")
	}

	if !utils.IsValidEmail(r.Email) {
		errors["email"] = append(errors["email"], "Invalid email format")
	}

	if !utils.IsValidPassword(r.Password) {
		errors["password"] = append(errors["password"], "Password must be at least 8 characters, with uppercase, number, and special character")
	}

	if r.Password != r.PasswordConfirm {
		errors["password_confirm"] = append(errors["password_confirm"], "Passwords do not match")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
