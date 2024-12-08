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

func RegisterUser(username, name, email, phone, password, confirmPassword, roleId string) error {
	if password != confirmPassword {
		return errors.New("password dan confirm password tidak cocok")
	}

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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
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

func UpdateUser(userID, email, username, name, phone string) error {

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if email != "" && email != user.Email && repositories.IsEmailExist(email) {
		return errors.New("email is already registered")
	}

	if username != "" && username != user.Username && repositories.IsUsernameExist(username) {
		return errors.New("username is already registered")
	}

	if phone != "" && phone != user.Phone && repositories.IsPhoneExist(phone) {
		return errors.New("phone number is already registered")
	}

	if email != "" {
		user.Email = email
	}
	if username != "" {
		user.Username = username
	}
	if name != "" {
		user.Name = name
	}
	if phone != "" {
		user.Phone = phone
	}

	err = repositories.UpdateUser(&user)
	if err != nil {
		return errors.New("failed to update user")
	}

	return nil
}

func UpdatePassword(userID, oldPassword, newPassword string) error {

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	err = repositories.UpdateUserPassword(userID, string(hashedPassword))
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
