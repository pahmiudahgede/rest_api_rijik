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
	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(req *dto.RegisterRequest) error
	VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) AuthService {
	return &authService{userRepo, roleRepo}
}

const otpCooldown = 30

func (s *authService) RegisterUser(req *dto.RegisterRequest) error {

	user, err := s.userRepo.GetUserByPhone(req.Phone)
	if err == nil && user != nil {
		return errors.New("phone number already registered")
	}

	lastOtpSent, err := utils.GetStringData("otp_sent:" + req.Phone)
	if err == nil && lastOtpSent != "" {
		lastSentTime, err := time.Parse(time.RFC3339, lastOtpSent)
		if err != nil {
			return errors.New("invalid OTP sent timestamp")
		}

		if time.Since(lastSentTime).Seconds() < otpCooldown {
			return errors.New("please wait before requesting another OTP")
		}
	}

	userID := uuid.New().String()

	user = &model.User{
		Phone:  req.Phone,
		RoleID: req.RoleID,
	}

	err = utils.SetJSONData("user:"+userID, user, 10*time.Minute)
	if err != nil {
		return err
	}

	err = utils.SetStringData("user_phone:"+req.Phone, userID, 10*time.Minute)
	if err != nil {
		return err
	}

	otp := generateOTP()

	err = config.SendWhatsAppMessage(req.Phone, fmt.Sprintf("Your OTP is: %s", otp))
	if err != nil {
		return err
	}

	err = utils.SetStringData("otp:"+req.Phone, otp, 10*time.Minute)
	if err != nil {
		return err
	}

	err = utils.SetStringData("otp_sent:"+req.Phone, time.Now().Format(time.RFC3339), 10*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error) {

	isLoggedIn, err := utils.GetStringData("user_logged_in:" + req.Phone)
	if err == nil && isLoggedIn == "true" {
		return nil, errors.New("you are already logged in")
	}

	storedOTP, err := utils.GetStringData("otp:" + req.Phone)
	if err != nil {
		return nil, err
	}

	if storedOTP == "" {
		return nil, errors.New("OTP expired or not found")
	}

	if storedOTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	userID, err := utils.GetStringData("user_phone:" + req.Phone)
	if err != nil || userID == "" {
		return nil, errors.New("user data not found in Redis")
	}

	userData, err := utils.GetJSONData("user:" + userID)
	if err != nil || userData == nil {
		return nil, errors.New("user data not found in Redis")
	}

	user := &model.User{
		Phone:  userData["phone"].(string),
		RoleID: userData["roleId"].(string),
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByID(createdUser.RoleID)
	if err != nil {
		return nil, err
	}

	token, err := generateJWTToken(createdUser.ID)
	if err != nil {
		return nil, err
	}

	err = utils.SetStringData("user_logged_in:"+req.Phone, "true", 0)
	if err != nil {
		return nil, err
	}

	return &dto.UserDataResponse{
		UserID:   createdUser.ID,
		UserRole: role.RoleName,
		Token:    token,
	}, nil
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}

func generateJWTToken(userID string) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Issuer:    userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := config.GetSecretKey()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
