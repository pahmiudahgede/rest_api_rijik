package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"math/rand"

	"rijig/config"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
)

type AuthService interface {
	RegisterUser(request *dto.RegisterRequest) error
	VerifyOTP(request *dto.VerifyOTPRequest) error
}

type authService struct {
	userRepo  repositories.UserRepository
	roleRepo  repositories.RoleRepository
	redisRepo repositories.RedisRepository
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, redisRepo repositories.RedisRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		redisRepo: redisRepo,
	}
}

func (s *authService) RegisterUser(request *dto.RegisterRequest) error {

	if request.RoleID == "" {
		return fmt.Errorf("role_id cannot be empty")
	}

	role, err := s.roleRepo.FindByID(request.RoleID)
	if err != nil {
		return fmt.Errorf("role not found: %v", err)
	}
	if role == nil {
		return fmt.Errorf("role with ID %s not found", request.RoleID)
	}

	existingUser, err := s.userRepo.FindByPhone(request.Phone)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("phone number already registered")
	}

	temporaryData := &model.User{
		Phone:  request.Phone,
		RoleID: request.RoleID,
	}

	err = s.redisRepo.StoreData(request.Phone, temporaryData, 1*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to store registration data in Redis: %v", err)
	}

	otp := generateOTP()
	err = s.redisRepo.StoreData("otp:"+request.Phone, otp, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %v", err)
	}

	err = config.SendWhatsAppMessage(request.Phone, fmt.Sprintf("Your OTP is: %s", otp))
	if err != nil {
		return fmt.Errorf("failed to send OTP via WhatsApp: %v", err)
	}

	log.Printf("OTP sent to phone number: %s", request.Phone)
	return nil
}

func (s *authService) VerifyOTP(request *dto.VerifyOTPRequest) error {

	storedOTP, err := s.redisRepo.GetData("otp:" + request.Phone)
	if err != nil {
		return fmt.Errorf("failed to retrieve OTP from Redis: %v", err)
	}
	if storedOTP != request.OTP {
		return fmt.Errorf("invalid OTP")
	}

	temporaryData, err := s.redisRepo.GetData(request.Phone)
	if err != nil {
		return fmt.Errorf("failed to get registration data from Redis: %v", err)
	}
	if temporaryData == "" {
		return fmt.Errorf("no registration data found for phone: %s", request.Phone)
	}

	temporaryDataStr, ok := temporaryData.(string)
	if !ok {
		return fmt.Errorf("failed to assert data to string")
	}

	temporaryDataBytes := []byte(temporaryDataStr)

	var user model.User
	err = json.Unmarshal(temporaryDataBytes, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal registration data: %v", err)
	}

	_, err = s.userRepo.SaveUser(&user)
	if err != nil {
		return fmt.Errorf("failed to save user to database: %v", err)
	}

	err = s.redisRepo.DeleteData(request.Phone)
	if err != nil {
		return fmt.Errorf("failed to delete registration data from Redis: %v", err)
	}

	err = s.redisRepo.DeleteData("otp:" + request.Phone)
	if err != nil {
		return fmt.Errorf("failed to delete OTP from Redis: %v", err)
	}

	return nil
}

func generateOTP() string {

	return fmt.Sprintf("%06d", RandomInt(100000, 999999))
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
