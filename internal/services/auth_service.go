package services

import (
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
	SecretKey string
}

func NewUserService(userRepo repositories.UserRepository, secretKey string) UserService {
	return &userService{UserRepo: userRepo, SecretKey: secretKey}
}

func (s *userService) Login(credentials dto.LoginDTO) (*dto.UserResponseWithToken, error) {

	user, err := s.UserRepo.FindByEmailOrUsernameOrPhone(credentials.Identifier)
	if err != nil {

		return nil, fmt.Errorf("user not found")
	}

	if !CheckPasswordHash(credentials.Password, user.Password) {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	err = utils.SetData(credentials.Identifier, token, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponseWithToken{
		UserID: user.ID,
		Token:  token,
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
	}

	err = s.UserRepo.Create(&newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &newUser, nil
}
