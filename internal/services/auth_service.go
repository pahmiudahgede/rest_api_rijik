package services

import (
	"errors"
	"fmt"
	"math/rand"
	"rijig/config"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const otpCooldown = 30 * time.Second

type AuthService interface {
	RegisterOrLogin(req *dto.RegisterRequest) error
	VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error)
	Logout(userID, phone string) error
}

type authService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) AuthService {
	return &authService{userRepo, roleRepo}
}

func (s *authService) RegisterOrLogin(req *dto.RegisterRequest) error {

	if err := s.checkOTPRequestCooldown(req.Phone); err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByPhoneAndRole(req.Phone, req.RoleID)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if user != nil {
		return s.sendOTP(req.Phone)
	}

	user = &model.User{
		Phone:  req.Phone,
		RoleID: req.RoleID,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}

	if err := s.saveUserToRedis(createdUser.ID, createdUser, req.Phone); err != nil {
		return err
	}

	return s.sendOTP(req.Phone)
}

func (s *authService) checkOTPRequestCooldown(phone string) error {
	otpSentTime, err := utils.GetStringData("otp_sent:" + phone)
	if err != nil || otpSentTime == "" {
		return nil
	}
	lastSent, _ := time.Parse(time.RFC3339, otpSentTime)
	if time.Since(lastSent) < otpCooldown {
		return errors.New("please wait before requesting a new OTP")
	}
	return nil
}

func (s *authService) sendOTP(phone string) error {
	otp := generateOTP()
	if err := config.SendWhatsAppMessage(phone, fmt.Sprintf("Your OTP is: %s", otp)); err != nil {
		return err
	}

	if err := utils.SetStringData("otp:"+phone, otp, 10*time.Minute); err != nil {
		return err
	}
	return utils.SetStringData("otp_sent:"+phone, time.Now().Format(time.RFC3339), 10*time.Minute)
}

func (s *authService) VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error) {

	storedOTP, err := utils.GetStringData("otp:" + req.Phone)
	if err != nil || storedOTP == "" {
		return nil, errors.New("OTP expired or not found")
	}

	if storedOTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	if err := utils.DeleteData("otp:" + req.Phone); err != nil {
		return nil, fmt.Errorf("failed to remove OTP from Redis: %w", err)
	}

	existingUser, err := s.userRepo.GetUserByPhoneAndRole(req.Phone, req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	var user *model.User
	if existingUser != nil {
		user = existingUser
	} else {

		user = &model.User{
			Phone:  req.Phone,
			RoleID: req.RoleID,
		}
		createdUser, err := s.userRepo.CreateUser(user)
		if err != nil {
			return nil, err
		}
		user = createdUser
	}

	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	if err := s.saveSessionData(user.ID, user.RoleID, role.RoleName, token); err != nil {
		return nil, err
	}

	return &dto.UserDataResponse{
		UserID:   user.ID,
		UserRole: role.RoleName,
		Token:    token,
	}, nil
}

func (s *authService) saveUserToRedis(userID string, user *model.User, phone string) error {
	if err := utils.SetJSONData("user:"+userID, user, 10*time.Minute); err != nil {
		return fmt.Errorf("failed to store user data in Redis: %w", err)
	}

	if err := utils.SetStringData("user_phone:"+userID, phone, 10*time.Minute); err != nil {
		return fmt.Errorf("failed to store user phone in Redis: %w", err)
	}

	return nil
}

func (s *authService) generateJWTToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := config.GetSecretKey()

	return token.SignedString([]byte(secretKey))
}

func (s *authService) saveSessionData(userID string, roleID string, roleName string, token string) error {
	sessionKey := fmt.Sprintf("session:%s", userID)
	sessionData := map[string]interface{}{
		"userID":   userID,
		"roleID":   roleID,
		"roleName": roleName,
	}

	if err := utils.SetJSONData(sessionKey, sessionData, 24*time.Hour); err != nil {
		return fmt.Errorf("failed to set session data: %w", err)
	}

	if err := utils.SetStringData("session_token:"+userID, token, 24*time.Hour); err != nil {
		return fmt.Errorf("failed to set session token: %w", err)
	}

	return nil
}

func (s *authService) Logout(userID, phone string) error {
	keys := []string{
		"session:" + userID,
		"session_token:" + userID,
		"user_logged_in:" + userID,
		"user:" + userID,
		"user_phone:" + userID,
		"otp_sent:" + phone,
	}

	for _, key := range keys {
		if err := utils.DeleteData(key); err != nil {
			return fmt.Errorf("failed to delete key %s from Redis: %w", key, err)
		}
	}

	return nil
}

func generateOTP() string {
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%04d", randGenerator.Intn(10000))
}
