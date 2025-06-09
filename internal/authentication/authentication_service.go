package authentication

import (
	"context"
	"fmt"
	"strings"
	"time"

	"rijig/internal/role"
	"rijig/model"
	"rijig/utils"
)

type AuthenticationService interface {
	LoginAdmin(ctx context.Context, req *LoginAdminRequest) (*AuthResponse, error)
	RegisterAdmin(ctx context.Context, req *RegisterAdminRequest) error

	SendRegistrationOTP(ctx context.Context, req *LoginorRegistRequest) (*OTPResponse, error)
	VerifyRegistrationOTP(ctx context.Context, req *VerifyOtpRequest) (*AuthResponse, error)

	SendLoginOTP(ctx context.Context, req *LoginorRegistRequest) (*OTPResponse, error)
	VerifyLoginOTP(ctx context.Context, req *VerifyOtpRequest) (*AuthResponse, error)

	LogoutAuthentication(ctx context.Context, userID, deviceID string) error
}

type authenticationService struct {
	authRepo AuthenticationRepository
	roleRepo role.RoleRepository
}

func NewAuthenticationService(authRepo AuthenticationRepository, roleRepo role.RoleRepository) AuthenticationService {
	return &authenticationService{authRepo, roleRepo}
}

func normalizeRoleName(roleName string) string {
	switch strings.ToLower(roleName) {
	case "administrator", "admin":
		return utils.RoleAdministrator
	case "pengelola":
		return utils.RolePengelola
	case "pengepul":
		return utils.RolePengepul
	case "masyarakat":
		return utils.RoleMasyarakat
	default:
		return strings.ToLower(roleName)
	}
}

func (s *authenticationService) LoginAdmin(ctx context.Context, req *LoginAdminRequest) (*AuthResponse, error) {
	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.Role == nil || user.Role.RoleName != "administrator" {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.RegistrationStatus != "completed" {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !utils.CompareHashAndPlainText(user.Password, req.Password) {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	token, err := utils.GenerateTokenPair(user.ID, user.Role.RoleName, req.DeviceID, user.RegistrationStatus, int(user.RegistrationProgress))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Message:            "login berhasil",
		AccessToken:        token.AccessToken,
		RefreshToken:       token.RefreshToken,
		RegistrationStatus: user.RegistrationStatus,
		SessionID:          token.SessionID,
	}, nil
}

func (s *authenticationService) RegisterAdmin(ctx context.Context, req *RegisterAdminRequest) error {

	existingUser, _ := s.authRepo.FindUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return fmt.Errorf("email already in use")
	}

	hashedPassword, err := utils.HashingPlainText(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	role, err := s.roleRepo.FindRoleByName(ctx, "administrator")
	if err != nil {
		return fmt.Errorf("role name not found: %w", err)
	}

	user := &model.User{
		Name:               req.Name,
		Phone:              req.Phone,
		Email:              req.Email,
		Gender:             req.Gender,
		Dateofbirth:        req.DateOfBirth,
		Placeofbirth:       req.PlaceOfBirth,
		Password:           hashedPassword,
		RoleID:             role.ID,
		RegistrationStatus: "completed",
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *authenticationService) SendRegistrationOTP(ctx context.Context, req *LoginorRegistRequest) (*OTPResponse, error) {

	normalizedRole := strings.ToLower(req.RoleName)

	existingUser, err := s.authRepo.FindUserByPhoneAndRole(ctx, req.Phone, normalizedRole)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("nomor telepon dengan role %s sudah terdaftar", req.RoleName)
	}

	roleData, err := s.roleRepo.FindRoleByName(ctx, normalizedRole)
	if err != nil {
		return nil, fmt.Errorf("role tidak valid: %v", err)
	}

	rateLimitKey := fmt.Sprintf("otp_limit:%s", req.Phone)
	if isRateLimited(rateLimitKey, 3, 5*time.Minute) {
		return nil, fmt.Errorf("terlalu banyak permintaan OTP, coba lagi dalam 5 menit")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("gagal generate OTP: %v", err)
	}

	otpKey := fmt.Sprintf("otp:%s:register", req.Phone)
	otpData := OTPData{
		Phone:     req.Phone,
		OTP:       otp,
		Role:      normalizedRole,
		RoleID:    roleData.ID,
		Type:      "register",
		Attempts:  0,
		ExpiresAt: time.Now().Add(90 * time.Second),
	}

	err = utils.SetCacheWithTTL(otpKey, otpData, 90*time.Second)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan OTP: %v", err)
	}

	err = sendOTP(req.Phone, otp)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim OTP: %v", err)
	}

	return &OTPResponse{
		Message:   "OTP berhasil dikirim",
		ExpiresIn: 90,
		Phone:     maskPhoneNumber(req.Phone),
	}, nil
}

func (s *authenticationService) VerifyRegistrationOTP(ctx context.Context, req *VerifyOtpRequest) (*AuthResponse, error) {

	otpKey := fmt.Sprintf("otp:%s:register", req.Phone)
	var otpData OTPData
	err := utils.GetCache(otpKey, &otpData)
	if err != nil {
		return nil, fmt.Errorf("OTP tidak ditemukan atau sudah kadaluarsa")
	}

	if otpData.Attempts >= 3 {
		utils.DeleteCache(otpKey)
		return nil, fmt.Errorf("terlalu banyak percobaan, silakan minta OTP baru")
	}

	if otpData.OTP != req.Otp {
		otpData.Attempts++
		utils.SetCacheWithTTL(otpKey, otpData, time.Until(otpData.ExpiresAt))
		return nil, fmt.Errorf("kode OTP salah")
	}

	if otpData.Role != req.RoleName {
		return nil, fmt.Errorf("role tidak sesuai")
	}

	normalizedRole := strings.ToLower(req.RoleName)

	user := &model.User{
		Phone:                req.Phone,
		PhoneVerified:        true,
		RoleID:               otpData.RoleID,
		RegistrationStatus:   utils.RegStatusIncomplete,
		RegistrationProgress: utils.ProgressOTPVerified,
		Name:                 "",
		Gender:               "",
		Dateofbirth:          "",
		Placeofbirth:         "",
	}

	err = s.authRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat user: %v", err)
	}

	if user.ID == "" {
		return nil, fmt.Errorf("gagal mendapatkan user ID setelah registrasi")
	}

	utils.DeleteCache(otpKey)

	tokenResponse, err := utils.GenerateTokenPair(
		user.ID,
		normalizedRole,
		req.DeviceID,
		user.RegistrationStatus,
		int(user.RegistrationProgress),
	)

	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %v", err)
	}

	nextStep := utils.GetNextRegistrationStep(
		normalizedRole,
		int(user.RegistrationProgress),
		user.RegistrationStatus,
	)

	return &AuthResponse{
		Message:            "Registrasi berhasil",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: user.RegistrationStatus,
		NextStep:           nextStep,
		SessionID:          tokenResponse.SessionID,
	}, nil
}

func (s *authenticationService) SendLoginOTP(ctx context.Context, req *LoginorRegistRequest) (*OTPResponse, error) {

	user, err := s.authRepo.FindUserByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("nomor telepon tidak terdaftar")
	}

	if !user.PhoneVerified {
		return nil, fmt.Errorf("nomor telepon belum diverifikasi")
	}

	rateLimitKey := fmt.Sprintf("otp_limit:%s", req.Phone)
	if isRateLimited(rateLimitKey, 3, 5*time.Minute) {
		return nil, fmt.Errorf("terlalu banyak permintaan OTP, coba lagi dalam 5 menit")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("gagal generate OTP: %v", err)
	}

	otpKey := fmt.Sprintf("otp:%s:login", req.Phone)
	otpData := OTPData{
		Phone:    req.Phone,
		OTP:      otp,
		UserID:   user.ID,
		Role:     user.Role.RoleName,
		Type:     "login",
		Attempts: 0,
	}

	err = utils.SetCacheWithTTL(otpKey, otpData, 1*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan OTP: %v", err)
	}

	err = sendOTP(req.Phone, otp)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim OTP: %v", err)
	}

	return &OTPResponse{
		Message:   "OTP berhasil dikirim",
		ExpiresIn: 300,
		Phone:     maskPhoneNumber(req.Phone),
	}, nil
}

func (s *authenticationService) VerifyLoginOTP(ctx context.Context, req *VerifyOtpRequest) (*AuthResponse, error) {

	otpKey := fmt.Sprintf("otp:%s:login", req.Phone)
	var otpData OTPData
	err := utils.GetCache(otpKey, &otpData)
	if err != nil {
		return nil, fmt.Errorf("OTP tidak ditemukan atau sudah kadaluarsa")
	}

	if otpData.Attempts >= 3 {
		utils.DeleteCache(otpKey)
		return nil, fmt.Errorf("terlalu banyak percobaan, silakan minta OTP baru")
	}

	if otpData.OTP != req.Otp {
		otpData.Attempts++
		utils.SetCache(otpKey, otpData, time.Until(otpData.ExpiresAt))
		return nil, fmt.Errorf("kode OTP salah")
	}

	user, err := s.authRepo.FindUserByID(ctx, otpData.UserID)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}

	utils.DeleteCache(otpKey)

	tokenResponse, err := utils.GenerateTokenPair(
		user.ID,
		user.Role.RoleName,
		req.DeviceID,
		"pin_verification_required",
		int(user.RegistrationProgress),
	)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %v", err)
	}

	return &AuthResponse{
		Message:            "OTP berhasil diverifikasi, silakan masukkan PIN",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		User:               convertUserToResponse(user),
		RegistrationStatus: user.RegistrationStatus,
		NextStep:           "Masukkan PIN",
		SessionID:          tokenResponse.SessionID,
	}, nil
}

func (s *authenticationService) LogoutAuthentication(ctx context.Context, userID, deviceID string) error {
	if err := utils.RevokeRefreshToken(userID, deviceID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

func maskPhoneNumber(phone string) string {
	if len(phone) < 4 {
		return phone
	}
	return phone[:4] + strings.Repeat("*", len(phone)-8) + phone[len(phone)-4:]
}

func isRateLimited(key string, maxAttempts int, duration time.Duration) bool {
	var count int
	err := utils.GetCache(key, &count)
	if err != nil {
		count = 0
	}

	if count >= maxAttempts {
		return true
	}

	count++
	utils.SetCache(key, count, duration)
	return false
}

func sendOTP(phone, otp string) error {

	fmt.Printf("Sending OTP %s to %s\n", otp, phone)
	return nil
}

func convertUserToResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:                   user.ID,
		Name:                 user.Name,
		Phone:                user.Phone,
		Email:                user.Email,
		Role:                 user.Role.RoleName,
		RegistrationStatus:   user.RegistrationStatus,
		RegistrationProgress: user.RegistrationProgress,
		PhoneVerified:        user.PhoneVerified,
		Avatar:               user.Avatar,
	}
}

func IsRegistrationComplete(role string, progress int) bool {
	switch role {
	case "masyarakat":
		return progress >= 1
	case "pengepul":
		return progress >= 2
	case "pengelola":
		return progress >= 3
	}
	return false
}
