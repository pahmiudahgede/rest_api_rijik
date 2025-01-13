package repositories

import (
	"errors"
	"fmt"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func IsEmailExist(email, roleId string) bool {
	var user domain.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		if user.RoleID == roleId {
			return true
		}
	}
	return false
}

func IsUsernameExist(username, roleId string) bool {
	var user domain.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err == nil {
		if user.RoleID == roleId {
			return true
		}
	}
	return false
}

func IsPhoneExist(phone, roleId string) bool {
	var user domain.User
	if err := config.DB.Where("phone = ?", phone).First(&user).Error; err == nil {
		if user.RoleID == roleId {
			return true
		}
	}
	return false
}

func CreateUser(username, name, email, phone, password, roleId string) error {
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

func GetUserByEmailUsernameOrPhone(identifier, roleId string) (domain.User, error) {
	var user domain.User
	err := config.DB.Where("email = ? OR username = ? OR phone = ?", identifier, identifier, identifier).First(&user).Error
	if err != nil {
		return user, errors.New("user not found")
	}

	if roleId != "" && user.RoleID != roleId {
		return user, errors.New("identifier found but role does not match")
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
		return errors.New("failed to update user")
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
