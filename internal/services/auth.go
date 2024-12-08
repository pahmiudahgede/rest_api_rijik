package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, name, email, phone, password, roleId string) error {

	if repositories.IsEmailExist(email) {
		return errors.New("email is already registered")
	}
	if repositories.IsUsernameExist(username) {
		return errors.New("username is already registered")
	}
	if repositories.IsPhoneExist(phone) {
		return errors.New("phone number is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = repositories.CreateUser(username, name, email, phone, string(hashedPassword), roleId)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(emailOrUsername, password string) (string, error) {
	if emailOrUsername == "" || password == "" {
		return "", errors.New("email/username and password must be provided")
	}

	user, err := repositories.GetUserByEmailOrUsername(emailOrUsername)
	if err != nil {
		return "", errors.New("invalid email/username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email/username or password")
	}

	token := generateJWT(user.ID)

	return token, nil
}

func generateJWT(userID string) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("API_KEY")))
	if err != nil {
		return ""
	}

	return t
}

func GetUserByID(userID string) (domain.User, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}
