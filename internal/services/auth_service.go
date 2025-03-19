package services

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
)

type AuthService interface {
	RegisterUser(request dto.RegisterRequest) (*model.User, error)
	VerifyOTP(phone, otp string) error
	GetUserByPhone(phone string) (*model.User, error)
	GenerateJWT(user *model.User) (string, error)
}

type authService struct {
	UserRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{UserRepo: userRepo}
}

func (s *authService) RegisterUser(request dto.RegisterRequest) (*model.User, error) {

	user, err := s.UserRepo.FindByPhone(request.Phone)
	if err == nil && user != nil {
		return nil, fmt.Errorf("user with phone %s already exists", request.Phone)
	}

	user = &model.User{
		Phone:         request.Phone,
		RoleID:        request.RoleID,
		EmailVerified: false,
	}

	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	_, err = s.SendOTP(request.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	return user, nil
}

func (s *authService) GetUserByPhone(phone string) (*model.User, error) {
	user, err := s.UserRepo.FindByPhone(phone)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user by phone: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *authService) SendOTP(phone string) (string, error) {
	otpCode := generateOTP()

	message := fmt.Sprintf("Your OTP code is: %s", otpCode)
	err := config.SendWhatsAppMessage(phone, message)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP via WhatsApp: %v", err)
	}

	expirationTime := 5 * time.Minute
	err = config.RedisClient.Set(config.Ctx, phone, otpCode, expirationTime).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store OTP in Redis: %v", err)
	}

	return otpCode, nil
}

func (s *authService) VerifyOTP(phone, otp string) error {

	otpRecord, err := config.RedisClient.Get(config.Ctx, phone).Result()
	if err == redis.Nil {

		return fmt.Errorf("OTP not found or expired")
	} else if err != nil {

		return fmt.Errorf("failed to retrieve OTP from Redis: %v", err)
	}

	if otp != otpRecord {
		return fmt.Errorf("invalid OTP")
	}

	err = config.RedisClient.Del(config.Ctx, phone).Err()
	if err != nil {
		return fmt.Errorf("failed to delete OTP from Redis: %v", err)
	}

	return nil
}

func (s *authService) GenerateJWT(user *model.User) (string, error) {
	if user == nil || user.Role == nil {
		return "", fmt.Errorf("user or user role is nil, cannot generate token")
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role.RoleName,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := config.GetSecretKey()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %v", err)
	}

	return tokenString, nil
}

func generateOTP() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}
