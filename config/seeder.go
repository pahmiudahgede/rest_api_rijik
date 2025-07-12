package config

import (
	"log"
	"rijig/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func SeedDefaultRoles(db *gorm.DB) error {
	log.Println("Starting roles seeder...")

	defaultRoles := []model.Role{
		{
			ID:       "8d841890-3962-4ceb-a82d-182f2f127442",
			RoleName: "administrator",
		},
		{
			ID:       "3871ee3d-1ac1-4fd5-88e0-c7005dcbcb85",
			RoleName: "pengelola",
		},
		{
			ID:       "4c366bf6-9806-476a-ab4c-329c104de3be",
			RoleName: "pengepul",
		},
		{
			ID:       "39eebc88-a322-4c1f-b0c7-d3572429c8db",
			RoleName: "masyarakat",
		},
	}

	for _, role := range defaultRoles {
		var existingRole model.Role
		result := db.Where("id = ? OR role_name = ?", role.ID, role.RoleName).First(&existingRole)

		if result.Error == nil {
			log.Printf("Role '%s' already exists, skipping", role.RoleName)
			continue
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("Error creating role '%s': %v", role.RoleName, err)
			return err
		}
		log.Printf("Role '%s' created successfully with ID: %s", role.RoleName, role.ID)
	}

	log.Println("Roles seeder completed successfully!")
	return nil
}

func SeedDefaultUser(db *gorm.DB) error {
	log.Println("Starting user seeder...")

	var existingUser model.User
	result := db.Where("phone = ?", "6287874527342").First(&existingUser)
	if result.Error == nil {
		log.Println("Default user already exists, skipping seeder")
		return nil
	}

	hashedPassword, err := HashPassword("Pahmi12345,")
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	defaultUser := model.User{
		Name:                 "Fahmi Kurniawan",
		Gender:               "laki-laki",
		Dateofbirth:          "02-09-2004",
		Placeofbirth:         "Banyuwangi",
		Phone:                "6287874527342",
		Email:                "pahmilucu123@gmail.com",
		EmailVerified:        true,
		PhoneVerified:        true,
		Password:             hashedPassword,
		RoleID:               "8d841890-3962-4ceb-a82d-182f2f127442",
		RegistrationStatus:   "complete",
		RegistrationProgress: 0,
	}

	if err := db.Create(&defaultUser).Error; err != nil {
		log.Printf("Error creating default user: %v", err)
		return err
	}

	log.Printf("Default user created successfully with ID: %s", defaultUser.ID)
	return nil
}

func SeedDynamicRoles(db *gorm.DB) error {
	log.Println("Starting dynamic roles seeder...")

	roleNames := []string{"administrator", "pengelola", "pengepul", "masyarakat"}

	for _, roleName := range roleNames {
		var existingRole model.Role
		result := db.Where("role_name = ?", roleName).First(&existingRole)

		if result.Error == nil {
			log.Printf("Role '%s' already exists, skipping", roleName)
			continue
		}

		newRole := model.Role{
			RoleName: roleName,
		}

		if err := db.Create(&newRole).Error; err != nil {
			log.Printf("Error creating role '%s': %v", roleName, err)
			return err
		}
		log.Printf("Role '%s' created successfully with ID: %s", roleName, newRole.ID)
	}

	log.Println("Dynamic roles seeder completed successfully!")
	return nil
}

func GetRoleIDByName(db *gorm.DB, roleName string) (string, error) {
	var role model.Role
	if err := db.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return "", err
	}
	return role.ID, nil
}

func SeedDefaultUserWithDynamicRole(db *gorm.DB) error {
	log.Println("Starting user seeder with dynamic role...")

	var existingUser model.User
	result := db.Where("phone = ?", "6287874527342").First(&existingUser)
	if result.Error == nil {
		log.Println("Default user already exists, skipping seeder")
		return nil
	}

	adminRoleID, err := GetRoleIDByName(db, "administrator")
	if err != nil {
		log.Printf("Error getting administrator role ID: %v", err)
		return err
	}

	hashedPassword, err := HashPassword("Pahmi12345,")
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	defaultUser := model.User{
		Name:                 "Fahmi Kurniawan",
		Gender:               "laki-laki",
		Dateofbirth:          "02-09-2004",
		Placeofbirth:         "Banyuwangi",
		Phone:                "6287874527342",
		Email:                "pahmilucu123@gmail.com",
		EmailVerified:        true,
		PhoneVerified:        true,
		Password:             hashedPassword,
		RoleID:               adminRoleID,
		RegistrationStatus:   "complete",
		RegistrationProgress: 0,
	}

	if err := db.Create(&defaultUser).Error; err != nil {
		log.Printf("Error creating default user: %v", err)
		return err
	}

	log.Printf("Default user created successfully with ID: %s", defaultUser.ID)
	return nil
}

func RunSeeders(db *gorm.DB) error {
	log.Println("Starting database seeders...")

	if err := SeedDefaultRoles(db); err != nil {
		return err
	}

	if err := SeedDefaultUser(db); err != nil {
		return err
	}

	log.Println("Database seeders completed successfully!")
	return nil
}
