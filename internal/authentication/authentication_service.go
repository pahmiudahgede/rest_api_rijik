package authentication

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"rijig/internal/role"
	"rijig/model"
	"rijig/utils"
)

type AuthenticationService interface {
	GetRegistrationStatus(ctx context.Context, userID, deviceID string) (*AuthResponse, error)
	LoginAdmin(ctx context.Context, req *LoginAdminRequest) (*OTPAdminResponse, error)
	RegisterAdmin(ctx context.Context, req *RegisterAdminRequest) (*EmailVerificationResponse, error)

	VerifyAdminOTP(ctx context.Context, req *VerifyAdminOTPRequest) (*AuthResponse, error)
	ResendAdminOTP(ctx context.Context, req *ResendAdminOTPRequest) (*OTPAdminResponse, error)

	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ResetPasswordResponse, error)
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error

	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error
	ResendEmailVerification(ctx context.Context, req *ResendVerificationRequest) (*EmailVerificationResponse, error)

	SendRegistrationOTP(ctx context.Context, req *LoginorRegistRequest) (*OTPResponse, error)
	VerifyRegistrationOTP(ctx context.Context, req *VerifyOtpRequest) (*AuthResponse, error)

	SendLoginOTP(ctx context.Context, req *LoginorRegistRequest) (*OTPResponse, error)
	VerifyLoginOTP(ctx context.Context, req *VerifyOtpRequest) (*AuthResponse, error)

	LogoutAuthentication(ctx context.Context, userID, deviceID string) error
}

type authenticationService struct {
	authRepo     AuthenticationRepository
	roleRepo     role.RoleRepository
	emailService *utils.EmailService
}

func NewAuthenticationService(authRepo AuthenticationRepository, roleRepo role.RoleRepository) AuthenticationService {
	return &authenticationService{
		authRepo:     authRepo,
		roleRepo:     roleRepo,
		emailService: utils.NewEmailService(),
	}
}

// func normalizeRoleName(roleName string) string {
// 	switch strings.ToLower(roleName) {
// 	case "administrator", "admin":
// 		return utils.RoleAdministrator
// 	case "pengelola":
// 		return utils.RolePengelola
// 	case "pengepul":
// 		return utils.RolePengepul
// 	case "masyarakat":
// 		return utils.RoleMasyarakat
// 	default:
// 		return strings.ToLower(roleName)
// 	}
// }

type GetRegistrationStatusResponse struct {
	UserID               string `json:"userId"`
	RegistrationStatus   string `json:"registrationStatus"`
	RegistrationProgress int8   `json:"registrationProgress"`
	Name                 string `json:"name"`
	Phone                string `json:"phone"`
	Role                 string `json:"role"`
}

func (s *authenticationService) GetRegistrationStatus(ctx context.Context, userID, deviceID string) (*AuthResponse, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.Role.RoleName == "" {
		return nil, fmt.Errorf("user role not found")
	}

	if user.RegistrationStatus == utils.RegStatusPending {
		log.Printf("⏳ User %s (%s) registration is still pending approval", user.Name, user.Phone)

		return &AuthResponse{
			Message:            "Your registration is currently under review. Please wait for approval.",
			RegistrationStatus: user.RegistrationStatus,
			NextStep:           "wait_for_approval",
		}, nil
	}

	if user.RegistrationStatus == utils.RegStatusConfirmed || user.RegistrationStatus == utils.RegStatusRejected {
		tokenResponse, err := utils.GenerateTokenPair(
			user.ID,
			user.Role.RoleName,
			deviceID,
			user.RegistrationStatus,
			int(user.RegistrationProgress),
		)
		if err != nil {
			log.Printf("GenerateTokenPair error: %v", err)
			return nil, fmt.Errorf("failed to generate token: %v", err)
		}

		nextStep := utils.GetNextRegistrationStep(
			user.Role.RoleName,
			int(user.RegistrationProgress),
			user.RegistrationStatus,
		)

		var message string
		if user.RegistrationStatus == utils.RegStatusConfirmed {
			message = "Registration approved successfully"
			log.Printf("✅ User %s (%s) registration approved - generating tokens", user.Name, user.Phone)
		} else if user.RegistrationStatus == utils.RegStatusRejected {
			message = "Registration has been rejected"
			log.Printf("❌ User %s (%s) registration rejected - generating tokens for rejection flow", user.Name, user.Phone)
		}

		return &AuthResponse{
			Message:            message,
			AccessToken:        tokenResponse.AccessToken,
			RefreshToken:       tokenResponse.RefreshToken,
			TokenType:          string(tokenResponse.TokenType),
			ExpiresIn:          tokenResponse.ExpiresIn,
			RegistrationStatus: user.RegistrationStatus,
			NextStep:           nextStep,
			SessionID:          tokenResponse.SessionID,
		}, nil
	}

	return nil, fmt.Errorf("unsupported registration status: %s", user.RegistrationStatus)
}

func (s *authenticationService) LoginAdmin(ctx context.Context, req *LoginAdminRequest) (*OTPAdminResponse, error) {

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if user.Role == nil || user.Role.RoleName != "administrator" {
		return nil, fmt.Errorf("invalid credentials")
	}

	if user.RegistrationStatus != "completed" {
		return nil, fmt.Errorf("account not activated")
	}

	if !utils.CompareHashAndPlainText(user.Password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	if utils.IsOTPValid(req.Email) {
		remaining, _ := utils.GetOTPRemainingTime(req.Email)
		return &OTPAdminResponse{
			Message:       "OTP sudah dikirim sebelumnya",
			Email:         req.Email,
			ExpiresIn:     remaining,
			RemainingTime: formatDuration(remaining),
			CanResend:     false,
			MaxAttempts:   utils.MAX_OTP_ATTEMPTS,
		}, nil
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP")
	}

	if err := utils.StoreOTP(req.Email, otp); err != nil {
		return nil, fmt.Errorf("failed to store OTP")
	}

	if err := s.emailService.SendOTPEmail(req.Email, user.Name, otp); err != nil {
		log.Printf("Failed to send OTP email: %v", err)
		return nil, fmt.Errorf("failed to send OTP email")
	}

	return &OTPAdminResponse{
		Message:       "Kode OTP berhasil dikirim ke email Anda",
		Email:         req.Email,
		ExpiresIn:     utils.OTP_EXPIRY,
		RemainingTime: formatDuration(utils.OTP_EXPIRY),
		CanResend:     false,
		MaxAttempts:   utils.MAX_OTP_ATTEMPTS,
	}, nil
}

func (s *authenticationService) VerifyAdminOTP(ctx context.Context, req *VerifyAdminOTPRequest) (*AuthResponse, error) {

	if err := utils.ValidateOTP(req.Email, req.OTP); err != nil {
		return nil, err
	}

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.EmailVerified {
		user.EmailVerified = true
		if err := s.authRepo.UpdateUser(ctx, user); err != nil {
			log.Printf("Failed to update email verification status: %v", err)
		}
	}

	token, err := utils.GenerateTokenPair(
		user.ID,
		user.Role.RoleName,
		req.DeviceID,
		user.RegistrationStatus,
		int(user.RegistrationProgress),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Message:            "Login berhasil",
		AccessToken:        token.AccessToken,
		RefreshToken:       token.RefreshToken,
		RegistrationStatus: user.RegistrationStatus,
		SessionID:          token.SessionID,
	}, nil
}

func (s *authenticationService) ResendAdminOTP(ctx context.Context, req *ResendAdminOTPRequest) (*OTPAdminResponse, error) {

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	if user.Role == nil || user.Role.RoleName != "administrator" {
		return nil, fmt.Errorf("not authorized")
	}

	if utils.IsOTPValid(req.Email) {
		remaining, _ := utils.GetOTPRemainingTime(req.Email)
		return nil, fmt.Errorf("OTP masih berlaku. Tunggu %s untuk kirim ulang", formatDuration(remaining))
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP")
	}

	if err := utils.StoreOTP(req.Email, otp); err != nil {
		return nil, fmt.Errorf("failed to store OTP")
	}

	if err := s.emailService.SendOTPEmail(req.Email, user.Name, otp); err != nil {
		log.Printf("Failed to send OTP email: %v", err)
		return nil, fmt.Errorf("failed to send OTP email")
	}

	return &OTPAdminResponse{
		Message:       "Kode OTP baru berhasil dikirim",
		Email:         req.Email,
		ExpiresIn:     utils.OTP_EXPIRY,
		RemainingTime: formatDuration(utils.OTP_EXPIRY),
		CanResend:     false,
		MaxAttempts:   utils.MAX_OTP_ATTEMPTS,
	}, nil
}

func (s *authenticationService) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error {

	verificationData, err := utils.ValidateEmailVerificationToken(req.Email, req.Token)
	if err != nil {
		return err
	}

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.ID != verificationData.UserID {
		return fmt.Errorf("invalid verification token")
	}

	if user.EmailVerified {
		return fmt.Errorf("email sudah terverifikasi sebelumnya")
	}

	user.EmailVerified = true
	user.RegistrationStatus = utils.RegStatusComplete
	// user.RegistrationProgress = 3

	if err := s.authRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	if err := utils.MarkEmailVerificationTokenAsUsed(req.Email); err != nil {
		log.Printf("Failed to mark verification token as used: %v", err)
	}

	return nil
}

func (s *authenticationService) ResendEmailVerification(ctx context.Context, req *ResendVerificationRequest) (*EmailVerificationResponse, error) {

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	if user.Role == nil || user.Role.RoleName != "administrator" {
		return nil, fmt.Errorf("not authorized")
	}

	if user.EmailVerified {
		return nil, fmt.Errorf("email sudah terverifikasi")
	}

	if utils.IsEmailVerificationTokenValid(req.Email) {
		remaining, _ := utils.GetEmailVerificationTokenRemainingTime(req.Email)
		return &EmailVerificationResponse{
			Message:       "Email verifikasi sudah dikirim sebelumnya",
			Email:         req.Email,
			ExpiresIn:     remaining,
			RemainingTime: formatDuration(remaining),
		}, nil
	}

	token, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token")
	}

	if err := utils.StoreEmailVerificationToken(req.Email, user.ID, token); err != nil {
		return nil, fmt.Errorf("failed to store verification token")
	}

	if err := s.emailService.SendEmailVerificationEmail(req.Email, user.Name, token); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		return nil, fmt.Errorf("failed to send verification email")
	}

	return &EmailVerificationResponse{
		Message:       "Email verifikasi berhasil dikirim ulang",
		Email:         req.Email,
		ExpiresIn:     utils.EMAIL_VERIFICATION_TOKEN_EXPIRY,
		RemainingTime: formatDuration(utils.EMAIL_VERIFICATION_TOKEN_EXPIRY),
	}, nil
}

func (s *authenticationService) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ResetPasswordResponse, error) {

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {

		return &ResetPasswordResponse{
			Message:       "Jika email terdaftar, link reset password akan dikirim",
			Email:         req.Email,
			ExpiresIn:     utils.RESET_TOKEN_EXPIRY,
			RemainingTime: formatDuration(utils.RESET_TOKEN_EXPIRY),
		}, nil
	}

	if user.Role == nil || user.Role.RoleName != "administrator" {

		return &ResetPasswordResponse{
			Message:       "Jika email terdaftar, link reset password akan dikirim",
			Email:         req.Email,
			ExpiresIn:     utils.RESET_TOKEN_EXPIRY,
			RemainingTime: formatDuration(utils.RESET_TOKEN_EXPIRY),
		}, nil
	}

	if utils.IsResetTokenValid(req.Email) {
		remaining, _ := utils.GetResetTokenRemainingTime(req.Email)
		return &ResetPasswordResponse{
			Message:       "Link reset password sudah dikirim sebelumnya",
			Email:         req.Email,
			ExpiresIn:     remaining,
			RemainingTime: formatDuration(remaining),
		}, nil
	}

	token, err := utils.GenerateResetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token")
	}

	if err := utils.StoreResetToken(req.Email, user.ID, token); err != nil {
		return nil, fmt.Errorf("failed to store reset token")
	}

	if err := s.emailService.SendResetPasswordEmail(req.Email, user.Name, token); err != nil {
		log.Printf("Failed to send reset password email: %v", err)
		return nil, fmt.Errorf("failed to send reset password email")
	}

	return &ResetPasswordResponse{
		Message:       "Link reset password berhasil dikirim ke email Anda",
		Email:         req.Email,
		ExpiresIn:     utils.RESET_TOKEN_EXPIRY,
		RemainingTime: formatDuration(utils.RESET_TOKEN_EXPIRY),
	}, nil
}

func (s *authenticationService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {

	resetData, err := utils.ValidateResetToken(req.Email, req.Token)
	if err != nil {
		return err
	}

	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.ID != resetData.UserID {
		return fmt.Errorf("invalid reset token")
	}

	hashedPassword, err := utils.HashingPlainText(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	if err := s.authRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := utils.MarkResetTokenAsUsed(req.Email); err != nil {
		log.Printf("Failed to mark reset token as used: %v", err)
	}

	if err := utils.RevokeAllRefreshTokens(user.ID); err != nil {
		log.Printf("Failed to revoke refresh tokens: %v", err)
	}

	return nil
}

func (s *authenticationService) RegisterAdmin(ctx context.Context, req *RegisterAdminRequest) (*EmailVerificationResponse, error) {

	existingUser, _ := s.authRepo.FindUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already in use")
	}

	hashedPassword, err := utils.HashingPlainText(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role, err := s.roleRepo.FindRoleByName(ctx, "administrator")
	if err != nil {
		return nil, fmt.Errorf("role name not found: %w", err)
	}

	user := &model.User{
		Name:                 req.Name,
		Phone:                req.Phone,
		Email:                req.Email,
		Gender:               req.Gender,
		Dateofbirth:          req.DateOfBirth,
		Placeofbirth:         req.PlaceOfBirth,
		Password:             hashedPassword,
		RoleID:               role.ID,
		RegistrationStatus:   "pending_email_verification",
		RegistrationProgress: 1,
		EmailVerified:        false,
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token")
	}

	if err := utils.StoreEmailVerificationToken(req.Email, user.ID, token); err != nil {
		return nil, fmt.Errorf("failed to store verification token")
	}

	if err := s.emailService.SendEmailVerificationEmail(req.Email, user.Name, token); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		return nil, fmt.Errorf("failed to send verification email")
	}

	return &EmailVerificationResponse{
		Message:       "Admin berhasil didaftarkan. Silakan cek email untuk verifikasi",
		Email:         req.Email,
		ExpiresIn:     utils.EMAIL_VERIFICATION_TOKEN_EXPIRY,
		RemainingTime: formatDuration(utils.EMAIL_VERIFICATION_TOKEN_EXPIRY),
	}, nil
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
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

	normalizedRole := strings.ToLower(req.RoleName)

	user, err := s.authRepo.FindUserByPhoneAndRole(ctx, req.Phone, normalizedRole)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
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

	var message string
	if user.RegistrationStatus == utils.RegStatusComplete {
		message = "verif pin"
		nextStep = "verif_pin"
	} else {
		message = "otp berhasil diverifikasi"
	}

	return &AuthResponse{
		Message:            message,
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: user.RegistrationStatus,
		NextStep:           nextStep,
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

// func convertUserToResponse(user *model.User) *UserResponse {
// 	return &UserResponse{
// 		ID:                   user.ID,
// 		Name:                 user.Name,
// 		Phone:                user.Phone,
// 		Email:                user.Email,
// 		Role:                 user.Role.RoleName,
// 		RegistrationStatus:   user.RegistrationStatus,
// 		RegistrationProgress: user.RegistrationProgress,
// 		PhoneVerified:        user.PhoneVerified,
// 		Avatar:               user.Avatar,
// 	}
// }

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
