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

type UserService interface {
	Login(credentials dto.LoginDTO) (*dto.UserResponseWithToken, error)
	Register(user dto.RegisterDTO) (*model.User, error)
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
		return nil, errors.New("roleId is required")
	}

	user, err := s.UserRepo.FindByIdentifierAndRole(credentials.Identifier, credentials.RoleID)
	if err != nil {
		return nil, errors.New("akun dengan role tersebut belum terdaftar")
	}

	if !CheckPasswordHash(credentials.Password, user.Password) {
		return nil, errors.New("password yang anda masukkan salah")
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

func (s *userService) Register(user dto.RegisterDTO) (*model.User, error) {
	if user.Password != user.ConfirmPassword {
		return nil, fmt.Errorf("password and confirm password do not match")
	}

	if user.RoleID == "" {
		return nil, fmt.Errorf("roleId is required")
	}

	role, err := s.RoleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("invalid roleId")
	}

	existingUser, _ := s.UserRepo.FindByUsername(user.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username is already taken")
	}

	existingPhone, _ := s.UserRepo.FindByPhoneAndRole(user.Phone, user.RoleID)
	if existingPhone != nil {
		return nil, fmt.Errorf("phone number is already used for this role")
	}

	existingEmail, _ := s.UserRepo.FindByEmailAndRole(user.Email, user.RoleID)
	if existingEmail != nil {
		return nil, fmt.Errorf("email is already used for this role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
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
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	newUser.Role = *role

	return &newUser, nil
}
