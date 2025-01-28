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