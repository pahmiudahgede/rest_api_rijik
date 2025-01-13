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

	if repositories.IsEmailExist(email, roleId) {
		return errors.New("email is already registered with the same role")
	}

	if repositories.IsUsernameExist(username, roleId) {
		return errors.New("username is already registered with the same role")
	}

	if repositories.IsPhoneExist(phone, roleId) {
		return errors.New("phone number is already registered with the same role")
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

func LoginUser(identifier, password string) (string, error) {
	if identifier == "" || password == "" {
		return "", errors.New("email/username/phone and password must be provided")
	}

	const roleId = ""

	user, err := repositories.GetUserByEmailUsernameOrPhone(identifier, roleId)
	if err != nil {
		return "", errors.New("invalid email/username/phone or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email/username/phone or password")
	}

	token := generateJWT(user.ID, user.RoleID)
	return token, nil
}

func generateJWT(userID, role string) string {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
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

	if email != "" && email != user.Email {
		if repositories.IsEmailExist(email, user.RoleID) {
			return errors.New("email is already registered with the same role")
		}
		user.Email = email
	}

	if username != "" && username != user.Username {
		if repositories.IsUsernameExist(username, user.RoleID) {
			return errors.New("username is already registered with the same role")
		}
		user.Username = username
	}

	if phone != "" && phone != user.Phone {
		if repositories.IsPhoneExist(phone, user.RoleID) {
			return errors.New("phone number is already registered with the same role")
		}
		user.Phone = phone
	}

	if name != "" {
		user.Name = name
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		return errors.New("new password cannot be the same as the old password")
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
