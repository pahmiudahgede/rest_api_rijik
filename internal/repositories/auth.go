package repositories

import (
	"errors"
	"fmt"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func IsEmailExist(email string) bool {
	var user domain.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return true
	}
	return false
}

func IsUsernameExist(username string) bool {
	var user domain.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err == nil {
		return true
	}
	return false
}

func IsPhoneExist(phone string) bool {
	var user domain.User
	if err := config.DB.Where("phone = ?", phone).First(&user).Error; err == nil {
		return true
	}
	return false
}

func CreateUser(username, name, email, phone, password, roleId string) error {

	if IsEmailExist(email) {
		return errors.New("email is already registered")
	}

	if IsUsernameExist(username) {
		return errors.New("username is already registered")
	}

	if IsPhoneExist(phone) {
		return errors.New("phone number is already registered")
	}

	user := domain.User{
		Username: username,
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
		RoleID:   roleId,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		return errors.New("failed to create user")
	}
	return nil
}

func GetUserByEmailOrUsername(emailOrUsername string) (domain.User, error) {
	var user domain.User
	if err := config.DB.Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}

func GetUserByID(userID string) (domain.User, error) {
	var user domain.User
	if err := config.DB.
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	fmt.Printf("User ID: %s, Role: %v\n", user.ID, user.Role)

	return user, nil
}

func UpdateUser(user *domain.User) error {
	if err := config.DB.Save(user).Error; err != nil {
		return errors.New("failed to save user")
	}
	return nil
}

func UpdateUserPassword(userID, newPassword string) error {
	var user domain.User

	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	user.Password = newPassword

	if err := config.DB.Save(&user).Error; err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
