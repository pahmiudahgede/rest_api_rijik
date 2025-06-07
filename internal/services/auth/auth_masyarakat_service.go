package service
/* 
import (
	"errors"
	"fmt"
	"rijig/config"
	"rijig/dto"
	"rijig/internal/repositories"
	repository "rijig/internal/repositories/auth"
	"rijig/model"
	"rijig/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMasyarakatService interface {
	RegisterOrLogin(req *dto.RegisterRequest) error
	VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error)
	Logout(userID, deviceID string) error
}

type authMasyarakatService struct {
	userRepo repository.AuthPengelolaRepository
	roleRepo repositories.RoleRepository
}

func NewAuthMasyarakatService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) AuthMasyarakatService {
	return &authMasyarakatService{userRepo, roleRepo}
}

func (s *authMasyarakatService) generateJWTToken(userID string, deviceID string) (string, error) {

	expirationTime := time.Now().Add(672 * time.Hour)

	claims := jwt.MapClaims{
		"sub":       userID,
		"exp":       expirationTime.Unix(),
		"device_id": deviceID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := config.GetSecretKey()

	return token.SignedString([]byte(secretKey))
}

func (s *authMasyarakatService) RegisterOrLogin(req *dto.RegisterRequest) error {
	if err := s.checkOTPRequestCooldown(req.Phone); err != nil {
		return err
	}
	return s.sendOTP(req.Phone)
}

func (s *authMasyarakatService) checkOTPRequestCooldown(phone string) error {
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

func (s *authMasyarakatService) sendOTP(phone string) error {
	otp := generateOTP()
	if err := config.SendWhatsAppMessage(phone, fmt.Sprintf("Your OTP is: %s", otp)); err != nil {
		return err
	}

	if err := utils.SetStringData("otp:"+phone, otp, 10*time.Minute); err != nil {
		return err
	}
	return utils.SetStringData("otp_sent:"+phone, time.Now().Format(time.RFC3339), 10*time.Minute)
}

func (s *authMasyarakatService) VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error) {

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
	if err := utils.DeleteData("otp_sent:" + req.Phone); err != nil {
		return nil, fmt.Errorf("failed to remove otp_sent from Redis: %w", err)
	}

	existingUser, err := s.userRepo.GetUserByPhoneAndRole(req.Phone, "0e5684e4-b214-4bd0-972f-3be80c4649a0")
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	var user *model.User
	if existingUser != nil {
		user = existingUser
	} else {

		user = &model.User{
			Phone:              req.Phone,
			RoleID:             "0e5684e4-b214-4bd0-972f-3be80c4649a0",
			PhoneVerified:      true,
			RegistrationStatus: "completed",
		}
		createdUser, err := s.userRepo.CreateUser(user)
		if err != nil {
			return nil, err
		}
		user = createdUser
	}

	token, err := s.generateJWTToken(user.ID, req.DeviceID)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	deviceID := req.DeviceID
	if err := s.saveSessionData(user.ID, deviceID, user.RoleID, role.RoleName, token); err != nil {
		return nil, err
	}

	return &dto.UserDataResponse{
		UserID:   user.ID,
		UserRole: role.RoleName,
		Token:    token,
	}, nil
}

func (s *authMasyarakatService) saveSessionData(userID string, deviceID string, roleID string, roleName string, token string) error {
	sessionKey := fmt.Sprintf("session:%s:%s", userID, deviceID)
	sessionData := map[string]interface{}{
		"userID":   userID,
		"roleID":   roleID,
		"roleName": roleName,
	}

	if err := utils.SetJSONData(sessionKey, sessionData, 24*time.Hour); err != nil {
		return fmt.Errorf("failed to set session data: %w", err)
	}

	if err := utils.SetStringData("session_token:"+userID+":"+deviceID, token, 24*time.Hour); err != nil {
		return fmt.Errorf("failed to set session token: %w", err)
	}

	return nil
}

func (s *authMasyarakatService) Logout(userID, deviceID string) error {

	err := utils.DeleteSessionData(userID, deviceID)
	if err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}

	return nil
}
 */