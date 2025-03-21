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

func (s *authService) RegisterUser(req *dto.RegisterRequest) error {

	userID := uuid.New().String()

	user := &model.User{
		Phone:  req.Phone,
		RoleID: req.RoleID,
	}

	err := utils.SetJSONData("user:"+userID, user, 10*time.Minute)
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

	return nil
}

func (s *authService) VerifyOTP(req *dto.VerifyOTPRequest) (*dto.UserDataResponse, error) {
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
