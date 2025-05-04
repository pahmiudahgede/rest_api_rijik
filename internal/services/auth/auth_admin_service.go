package service

import (
	"errors"
	"fmt"
	"log"
	"rijig/config"
	dto "rijig/dto/auth"
	"rijig/internal/repositories"
	repository "rijig/internal/repositories/auth"
	"rijig/model"
	"rijig/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrEmailTaken            = "email is already used"
	ErrPhoneTaken            = "phone number is already used"
	ErrInvalidPassword       = "password does not match"
	ErrRoleNotFound          = "role not found"
	ErrFailedToGenerateToken = "failed to generate token"
	ErrFailedToHashPassword  = "failed to hash password"
	ErrFailedToCreateUser    = "failed to create user"
	ErrIncorrectPassword     = "incorrect password"
	ErrAccountNotFound       = "account not found"
)

type AuthAdminService interface {
	RegisterAdmin(request *dto.RegisterAdminRequest) (*model.User, error)

	LoginAdmin(req *dto.LoginAdminRequest) (*dto.LoginResponse, error)
	LogoutAdmin(userID, deviceID string) error
}

type authAdminService struct {
	UserRepo  repository.AuthAdminRepository
	RoleRepo  repositories.RoleRepository
	SecretKey string
}

func NewAuthAdminService(userRepo repository.AuthAdminRepository, roleRepo repositories.RoleRepository, secretKey string) AuthAdminService {
	return &authAdminService{UserRepo: userRepo, RoleRepo: roleRepo, SecretKey: secretKey}
}

func (s *authAdminService) RegisterAdmin(request *dto.RegisterAdminRequest) (*model.User, error) {

	if existingUser, _ := s.UserRepo.FindByEmail(request.Email); existingUser != nil {
		return nil, errors.New(ErrEmailTaken)
	}

	if existingUser, _ := s.UserRepo.FindByPhone(request.Phone); existingUser != nil {
		return nil, errors.New(ErrPhoneTaken)
	}

	role, err := s.UserRepo.FindRoleByName("administrator")
	if err != nil {
		return nil, errors.New(ErrRoleNotFound)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, errors.New(ErrFailedToHashPassword)
	}

	user := &model.User{
		Name:         request.Name,
		Email:        request.Email,
		Phone:        request.Phone,
		Password:     string(hashedPassword),
		RoleID:       role.ID,
		Role:         role,
		Dateofbirth:  request.Dateofbirth,
		Placeofbirth: request.Placeofbirth,
		Gender:       request.Gender,
		RegistrationStatus: "completed",
	}

	createdUser, err := s.UserRepo.CreateUser(user)
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToCreateUser, err)
	}

	return createdUser, nil
}

func (s *authAdminService) LoginAdmin(req *dto.LoginAdminRequest) (*dto.LoginResponse, error) {

	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		log.Println("User not found:", err)
		return nil, errors.New(ErrAccountNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("Incorrect password:", err)
		return nil, errors.New(ErrIncorrectPassword)
	}

	existingUser, err := s.UserRepo.FindAdminByEmailandRoleid(req.Email, "42bdecce-f2ad-44ae-b3d6-883c1fbddaf7")
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	var adminUser *model.User
	if existingUser != nil {
		adminUser = existingUser
	} else {

		adminUser = &model.User{
			Email:              req.Email,
			RoleID:             "42bdecce-f2ad-44ae-b3d6-883c1fbddaf7",
		}
		createdUser, err := s.UserRepo.CreateUser(adminUser)
		if err != nil {
			return nil, err
		}
		adminUser = createdUser
	}

	token, err := s.generateJWTToken(adminUser.ID, req.Deviceid)
	if err != nil {
		return nil, err
	}

	role, err := s.RoleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	deviceID := req.Deviceid
	if err := s.saveSessionAdminData(user.ID, deviceID, user.RoleID, role.RoleName, token); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		UserID: user.ID,
		Role:   user.Role.RoleName,
		Token:  token,
	}, nil
}

func (s *authAdminService) saveSessionAdminData(userID string, deviceID string, roleID string, roleName string, token string) error {
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

func (s *authAdminService) generateJWTToken(userID string, deviceID string) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"sub":       userID,
		"exp":       expirationTime.Unix(),
		"device_id": deviceID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := config.GetSecretKey()

	return token.SignedString([]byte(secretKey))
}

func (s *authAdminService) LogoutAdmin(userID, deviceID string) error {

	err := utils.DeleteSessionData(userID, deviceID)
	if err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}

	return nil
}
