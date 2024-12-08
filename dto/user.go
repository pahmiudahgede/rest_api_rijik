package dto

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email harus diisi")
	}

	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("format email belum sesuai")
	}
	return nil
}

func ValidatePhone(phone string) error {
	if phone == "" {
		return errors.New("nomor telepon harus diisi")
	}

	phoneRegex := `^\+?[0-9]{10,15}$`
	re := regexp.MustCompile(phoneRegex)
	if !re.MatchString(phone) {
		return errors.New("nomor telepon tidak valid")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password harus diisi")
	}

	if len(password) < 8 {
		return errors.New("password minimal 8 karakter")
	}
	return nil
}

type RegisterUserInput struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	RoleId          string `json:"roleId"`
}

func (input *RegisterUserInput) Validate() error {

	if input.Username == "" {
		return errors.New("username harus diisi")
	}

	if input.Name == "" {
		return errors.New("nama harus diisi")
	}

	if err := ValidateEmail(input.Email); err != nil {
		return err
	}

	if err := ValidatePhone(input.Phone); err != nil {
		return err
	}

	if err := ValidatePassword(input.Password); err != nil {
		return err
	}

	if input.Password != input.ConfirmPassword {
		return errors.New("password dan confirm password tidak cocok")
	}

	if input.RoleId == "" {
		return errors.New("roleId harus diisi")
	}

	return nil
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (input *UpdatePasswordInput) Validate() error {

	if input.OldPassword == "" {
		return errors.New("old password must be provided")
	}

	if input.NewPassword == "" {
		return errors.New("new password must be provided")
	}

	if len(input.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}

	return nil
}

type UpdateUserInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

func (input *UpdateUserInput) Validate() error {

	if input.Email != "" {
		if err := ValidateEmail(input.Email); err != nil {
			return err
		}
	}

	if input.Username == "" {
		return errors.New("username harus diisi")
	}

	if input.Name == "" {
		return errors.New("name harus diisi")
	}

	if input.Phone != "" {
		if err := ValidatePhone(input.Phone); err != nil {
			return err
		}
	}

	return nil
}
