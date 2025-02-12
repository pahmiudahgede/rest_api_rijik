package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrUsernameTaken        = "username is already taken"
	ErrPhoneTaken           = "phone number is already used for this role"
	ErrEmailTaken           = "email is already used for this role"
	ErrInvalidRoleID        = "invalid roleId"
	ErrPasswordMismatch     = "password and confirm password do not match"
	ErrRoleIDRequired       = "roleId is required"
	ErrFailedToHashPassword = "failed to hash password"
	ErrFailedToCreateUser   = "failed to create user"
	ErrIncorrectPassword    = "incorrect password"
	ErrAccountNotFound      = "account not found"
)

type UserService interface {
	Login(credentials dto.LoginDTO) (*dto.UserResponseWithToken, error)
	Register(user dto.RegisterDTO) (*dto.UserResponseDTO, error)
}

type userService struct {
	UserRepo  repositories.UserRepository
	RoleRepo  repositories.RoleRepository
	SecretKey string
}

func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, secretKey string) UserService {
	return &userService{UserRepo: userRepo, RoleRepo: roleRepo, SecretKey: secretKey}
}

func (s *userService) Login(credentials dto.LoginDTO) (*dto.UserResponseWithToken, error) {
	if credentials.RoleID == "" {
		return nil, errors.New(ErrRoleIDRequired)
	}

	user, err := s.UserRepo.FindByIdentifierAndRole(credentials.Identifier, credentials.RoleID)
	if err != nil {
		return nil, errors.New(ErrAccountNotFound)
	}

	if !CheckPasswordHash(credentials.Password, user.Password) {
		return nil, errors.New(ErrIncorrectPassword)
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	sessionKey := fmt.Sprintf("session:%s", user.ID)
	sessionData := map[string]interface{}{
		"userID":   user.ID,
		"roleID":   user.RoleID,
		"roleName": user.Role.RoleName,
	}

	err = utils.SetJSONData(sessionKey, sessionData, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponseWithToken{
		RoleName: user.Role.RoleName,
		UserID:   user.ID,
		Token:    token,
	}, nil
}

func (s *userService) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *userService) Register(user dto.RegisterDTO) (*dto.UserResponseDTO, error) {

	if user.Password != user.ConfirmPassword {
		return nil, fmt.Errorf("%s", ErrPasswordMismatch)
	}

	if user.RoleID == "" {
		return nil, fmt.Errorf("%s", ErrRoleIDRequired)
	}

	role, err := s.RoleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrInvalidRoleID, err)
	}

	if existingUser, _ := s.UserRepo.FindByUsername(user.Username); existingUser != nil {
		return nil, fmt.Errorf("%s", ErrUsernameTaken)
	}

	if existingPhone, _ := s.UserRepo.FindByPhoneAndRole(user.Phone, user.RoleID); existingPhone != nil {
		return nil, fmt.Errorf("%s", ErrPhoneTaken)
	}

	if existingEmail, _ := s.UserRepo.FindByEmailAndRole(user.Email, user.RoleID); existingEmail != nil {
		return nil, fmt.Errorf("%s", ErrEmailTaken)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrFailedToHashPassword, err)
	}

	newUser := model.User{
		Username: user.Username,
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: string(hashedPassword),
		RoleID:   user.RoleID,
	}

	err = s.UserRepo.Create(&newUser)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrFailedToCreateUser, err)
	}

	userResponse := s.prepareUserResponse(newUser, role)

	return userResponse, nil
}

func (s *userService) prepareUserResponse(user model.User, role *model.Role) *dto.UserResponseDTO {

	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	return &dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		RoleName:      role.RoleName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
