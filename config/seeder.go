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
	log.Println("Starting default administrator user seeder...")

	var existingUser model.User
	result := db.Where("phone = ?", "6287874527342").First(&existingUser)
	if result.Error == nil {
		log.Println("Default administrator user already exists, skipping seeder")
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
		RoleID:               "8d841890-3962-4ceb-a82d-182f2f127442", // administrator
		RegistrationStatus:   "complete",
		RegistrationProgress: 0,
	}

	if err := db.Create(&defaultUser).Error; err != nil {
		log.Printf("Error creating default user: %v", err)
		return err
	}

	log.Printf("Default administrator user created successfully with ID: %s", defaultUser.ID)
	return nil
}

func SeedPengelolaUsers(db *gorm.DB) error {
	log.Println("Starting pengelola users seeder...")

	pengelolaUsers := []struct {
		User           model.User
		CompanyProfile model.CompanyProfile
	}{
		{
			User: model.User{
				Name:                 "Budi Santoso",
				Gender:               "laki-laki",
				Dateofbirth:          "15-03-1985",
				Placeofbirth:         "Jakarta",
				Phone:                "6281234567890",
				Email:                "budi.santoso@greenwaste.co.id",
				EmailVerified:        true,
				PhoneVerified:        true,
				Password:             "",                                     // akan di-hash
				RoleID:               "3871ee3d-1ac1-4fd5-88e0-c7005dcbcb85", // pengelola
				RegistrationStatus:   "complete",
				RegistrationProgress: 0,
			},
			CompanyProfile: model.CompanyProfile{
				CompanyName:        "Green Waste Management",
				CompanyAddress:     "Jl. Sudirman No. 123, Jakarta Pusat, DKI Jakarta",
				CompanyPhone:       "02112345678",
				CompanyEmail:       "info@greenwaste.co.id",
				CompanyLogo:        "https://example.com/logos/greenwaste.png",
				CompanyWebsite:     "https://greenwaste.co.id",
				TaxID:              "01.234.567.8-901.000",
				FoundedDate:        "12-05-2010",
				CompanyType:        "PT",
				CompanyDescription: "Perusahaan pengelola sampah terpercaya yang berfokus pada pengelolaan sampah ramah lingkungan dan berkelanjutan.",
			},
		},
		{
			User: model.User{
				Name:                 "Siti Nurhaliza",
				Gender:               "perempuan",
				Dateofbirth:          "08-11-1982",
				Placeofbirth:         "Surabaya",
				Phone:                "6282345678901",
				Email:                "siti.nurhaliza@ecowaste.co.id",
				EmailVerified:        true,
				PhoneVerified:        true,
				Password:             "",                                     // akan di-hash
				RoleID:               "3871ee3d-1ac1-4fd5-88e0-c7005dcbcb85", // pengelola
				RegistrationStatus:   "complete",
				RegistrationProgress: 0,
			},
			CompanyProfile: model.CompanyProfile{
				CompanyName:        "Eco Waste Solutions",
				CompanyAddress:     "Jl. Ahmad Yani No. 456, Surabaya, Jawa Timur",
				CompanyPhone:       "03123456789",
				CompanyEmail:       "contact@ecowaste.co.id",
				CompanyLogo:        "https://example.com/logos/ecowaste.png",
				CompanyWebsite:     "https://ecowaste.co.id",
				TaxID:              "02.345.678.9-012.000",
				FoundedDate:        "20-08-2015",
				CompanyType:        "CV",
				CompanyDescription: "Solusi pengelolaan sampah inovatif dengan teknologi terkini untuk menciptakan lingkungan yang lebih bersih dan sehat.",
			},
		},
	}

	for i, data := range pengelolaUsers {
		// Check if user already exists
		var existingUser model.User
		result := db.Where("phone = ?", data.User.Phone).First(&existingUser)
		if result.Error == nil {
			log.Printf("Pengelola user with phone '%s' already exists, skipping", data.User.Phone)
			continue
		}

		// Hash password
		hashedPassword, err := HashPassword("Pengelola123!")
		if err != nil {
			log.Printf("Error hashing password for pengelola user %d: %v", i+1, err)
			return err
		}
		data.User.Password = hashedPassword

		// Create user
		if err := db.Create(&data.User).Error; err != nil {
			log.Printf("Error creating pengelola user %d: %v", i+1, err)
			return err
		}

		// Create company profile
		data.CompanyProfile.UserID = data.User.ID
		if err := db.Create(&data.CompanyProfile).Error; err != nil {
			log.Printf("Error creating company profile for pengelola user %d: %v", i+1, err)
			return err
		}

		log.Printf("Pengelola user '%s' and company profile '%s' created successfully", data.User.Name, data.CompanyProfile.CompanyName)
	}

	log.Println("Pengelola users seeder completed successfully!")
	return nil
}

func SeedPengepulUsers(db *gorm.DB) error {
	log.Println("Starting pengepul users seeder...")

	pengepulUsers := []struct {
		User         model.User
		IdentityCard model.IdentityCard
	}{
		{
			User: model.User{
				Name:                 "Ahmad Rizki",
				Gender:               "laki-laki",
				Dateofbirth:          "25-07-1990",
				Placeofbirth:         "Bandung",
				Phone:                "6283456789012",
				Email:                "ahmad.rizki@gmail.com",
				EmailVerified:        true,
				PhoneVerified:        true,
				Password:             "",                                     // akan di-hash
				RoleID:               "4c366bf6-9806-476a-ab4c-329c104de3be", // pengepul
				RegistrationStatus:   "complete",
				RegistrationProgress: 0,
			},
			IdentityCard: model.IdentityCard{
				Identificationumber: "3273012507900001",
				Fullname:            "Ahmad Rizki Pratama",
				Placeofbirth:        "Bandung",
				Dateofbirth:         "25-07-1990",
				Gender:              "LAKI-LAKI",
				BloodType:           "O",
				Province:            "JAWA BARAT",
				District:            "KOTA BANDUNG",
				SubDistrict:         "COBLONG",
				Hamlet:              "CIPAGANTI",
				Village:             "CIPAGANTI",
				Neighbourhood:       "003/008",
				PostalCode:          "40132",
				Religion:            "ISLAM",
				Maritalstatus:       "BELUM KAWIN",
				Job:                 "WIRAUSAHA",
				Citizenship:         "WNI",
				Validuntil:          "SEUMUR HIDUP",
				Cardphoto:           "https://example.com/identities/ahmad_ktp.jpg",
			},
		},
		{
			User: model.User{
				Name:                 "Dewi Kusuma",
				Gender:               "perempuan",
				Dateofbirth:          "12-02-1988",
				Placeofbirth:         "Yogyakarta",
				Phone:                "6284567890123",
				Email:                "dewi.kusuma@gmail.com",
				EmailVerified:        true,
				PhoneVerified:        true,
				Password:             "",                                     // akan di-hash
				RoleID:               "4c366bf6-9806-476a-ab4c-329c104de3be", // pengepul
				RegistrationStatus:   "complete",
				RegistrationProgress: 0,
			},
			IdentityCard: model.IdentityCard{
				Identificationumber: "3404015202880002",
				Fullname:            "Dewi Kusuma Wardani",
				Placeofbirth:        "Yogyakarta",
				Dateofbirth:         "12-02-1988",
				Gender:              "PEREMPUAN",
				BloodType:           "A",
				Province:            "DI YOGYAKARTA",
				District:            "KOTA YOGYAKARTA",
				SubDistrict:         "MERGANGSAN",
				Hamlet:              "WIROBRAJAN",
				Village:             "WIROBRAJAN",
				Neighbourhood:       "005/010",
				PostalCode:          "55253",
				Religion:            "ISLAM",
				Maritalstatus:       "KAWIN",
				Job:                 "WIRAUSAHA",
				Citizenship:         "WNI",
				Validuntil:          "SEUMUR HIDUP",
				Cardphoto:           "https://example.com/identities/dewi_ktp.jpg",
			},
		},
	}

	for i, data := range pengepulUsers {
		// Check if user already exists
		var existingUser model.User
		result := db.Where("phone = ?", data.User.Phone).First(&existingUser)
		if result.Error == nil {
			log.Printf("Pengepul user with phone '%s' already exists, skipping", data.User.Phone)
			continue
		}

		// Hash password
		hashedPassword, err := HashPassword("Pengepul123!")
		if err != nil {
			log.Printf("Error hashing password for pengepul user %d: %v", i+1, err)
			return err
		}
		data.User.Password = hashedPassword

		// Create user
		if err := db.Create(&data.User).Error; err != nil {
			log.Printf("Error creating pengepul user %d: %v", i+1, err)
			return err
		}

		// Create identity card
		data.IdentityCard.UserID = data.User.ID
		if err := db.Create(&data.IdentityCard).Error; err != nil {
			log.Printf("Error creating identity card for pengepul user %d: %v", i+1, err)
			return err
		}

		log.Printf("Pengepul user '%s' and identity card created successfully", data.User.Name)
	}

	log.Println("Pengepul users seeder completed successfully!")
	return nil
}

func SeedMasyarakatUsers(db *gorm.DB) error {
	log.Println("Starting masyarakat users seeder...")

	masyarakatUsers := []model.User{
		{
			Name:                 "Andi Setiawan",
			Gender:               "laki-laki",
			Dateofbirth:          "18-06-1992",
			Placeofbirth:         "Malang",
			Phone:                "6285678901234",
			Email:                "andi.setiawan@gmail.com",
			EmailVerified:        true,
			PhoneVerified:        true,
			Password:             "",                                     // akan di-hash
			RoleID:               "39eebc88-a322-4c1f-b0c7-d3572429c8db", // masyarakat
			RegistrationStatus:   "complete",
			RegistrationProgress: 0,
		},
		{
			Name:                 "Maya Sari",
			Gender:               "perempuan",
			Dateofbirth:          "03-12-1995",
			Placeofbirth:         "Denpasar",
			Phone:                "6286789012345",
			Email:                "maya.sari@gmail.com",
			EmailVerified:        true,
			PhoneVerified:        true,
			Password:             "",                                     // akan di-hash
			RoleID:               "39eebc88-a322-4c1f-b0c7-d3572429c8db", // masyarakat
			RegistrationStatus:   "complete",
			RegistrationProgress: 0,
		},
		{
			Name:                 "Reza Pratama",
			Gender:               "laki-laki",
			Dateofbirth:          "22-09-1987",
			Placeofbirth:         "Medan",
			Phone:                "6287890123456",
			Email:                "reza.pratama@gmail.com",
			EmailVerified:        true,
			PhoneVerified:        true,
			Password:             "",                                     // akan di-hash
			RoleID:               "39eebc88-a322-4c1f-b0c7-d3572429c8db", // masyarakat
			RegistrationStatus:   "complete",
			RegistrationProgress: 0,
		},
	}

	for i, user := range masyarakatUsers {
		// Check if user already exists
		var existingUser model.User
		result := db.Where("phone = ?", user.Phone).First(&existingUser)
		if result.Error == nil {
			log.Printf("Masyarakat user with phone '%s' already exists, skipping", user.Phone)
			continue
		}

		// Hash password
		hashedPassword, err := HashPassword("Masyarakat123!")
		if err != nil {
			log.Printf("Error hashing password for masyarakat user %d: %v", i+1, err)
			return err
		}
		user.Password = hashedPassword

		// Create user
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating masyarakat user %d: %v", i+1, err)
			return err
		}

		log.Printf("Masyarakat user '%s' created successfully with ID: %s", user.Name, user.ID)
	}

	log.Println("Masyarakat users seeder completed successfully!")
	return nil
}

func GetRoleIDByName(db *gorm.DB, roleName string) (string, error) {
	var role model.Role
	if err := db.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return "", err
	}
	return role.ID, nil
}

func RunSeeders(db *gorm.DB) error {
	log.Println("Starting database seeders...")

	// Seed roles first
	if err := SeedDefaultRoles(db); err != nil {
		log.Printf("Error seeding roles: %v", err)
		return err
	}

	// Seed default administrator user
	if err := SeedDefaultUser(db); err != nil {
		log.Printf("Error seeding default user: %v", err)
		return err
	}

	// Seed pengelola users with company profiles
	if err := SeedPengelolaUsers(db); err != nil {
		log.Printf("Error seeding pengelola users: %v", err)
		return err
	}

	// Seed pengepul users with identity cards
	if err := SeedPengepulUsers(db); err != nil {
		log.Printf("Error seeding pengepul users: %v", err)
		return err
	}

	// Seed masyarakat users
	if err := SeedMasyarakatUsers(db); err != nil {
		log.Printf("Error seeding masyarakat users: %v", err)
		return err
	}

	log.Println("Database seeders completed successfully!")
	return nil
}
