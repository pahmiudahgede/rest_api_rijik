package repositories

import (
	"rijig/model"

	"gorm.io/gorm"
)

type AuthAdminRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindAdminByEmailandRoleid(email, roleId string) (*model.User, error)
	FindAdminByPhoneandRoleid(phone, roleId string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	FindByEmailOrPhone(identifier string) (*model.User, error)
	FindRoleByName(roleName string) (*model.Role, error)
	CreateUser(user *model.User) (*model.User, error)
}

type authAdminRepository struct {
	DB *gorm.DB
}

func NewAuthAdminRepository(db *gorm.DB) AuthAdminRepository {
	return &authAdminRepository{DB: db}
}

func (r *authAdminRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authAdminRepository) FindAdminByEmailandRoleid(email, roleId string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ? AND role_id = ?", email, roleId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authAdminRepository) FindAdminByPhoneandRoleid(phone, roleId string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("phone = ? AND role_id = ?", phone, roleId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authAdminRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authAdminRepository) FindByEmailOrPhone(identifier string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ? OR phone = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authAdminRepository) CreateUser(user *model.User) (*model.User, error) {
	err := r.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authAdminRepository) FindRoleByName(roleName string) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
