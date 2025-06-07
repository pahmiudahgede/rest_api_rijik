package identitycart

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rijig/internal/authentication"
	"rijig/model"
	"rijig/utils"
	"strings"
)

type IdentityCardService interface {
	CreateIdentityCard(ctx context.Context, userID string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*authentication.AuthResponse, error)
	GetIdentityCardByID(ctx context.Context, id string) (*ResponseIdentityCardDTO, error)
	GetIdentityCardsByUserID(ctx context.Context, userID string) ([]ResponseIdentityCardDTO, error)
	UpdateIdentityCard(ctx context.Context, userID string, id string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*ResponseIdentityCardDTO, error)
}

type identityCardService struct {
	identityRepo IdentityCardRepository
	authRepo     authentication.AuthenticationRepository
}

func NewIdentityCardService(identityRepo IdentityCardRepository, authRepo authentication.AuthenticationRepository) IdentityCardService {
	return &identityCardService{
		identityRepo: identityRepo,
		authRepo:     authRepo,
	}
}

func FormatResponseIdentityCard(identityCard *model.IdentityCard) (*ResponseIdentityCardDTO, error) {
	createdAt, _ := utils.FormatDateToIndonesianFormat(identityCard.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(identityCard.UpdatedAt)

	return &ResponseIdentityCardDTO{
		ID:                  identityCard.ID,
		UserID:              identityCard.UserID,
		Identificationumber: identityCard.Identificationumber,
		Placeofbirth:        identityCard.Placeofbirth,
		Dateofbirth:         identityCard.Dateofbirth,
		Gender:              identityCard.Gender,
		BloodType:           identityCard.BloodType,
		Province:            identityCard.Province,
		District:            identityCard.District,
		SubDistrict:         identityCard.SubDistrict,
		Hamlet:              identityCard.Hamlet,
		Village:             identityCard.Village,
		Neighbourhood:       identityCard.Neighbourhood,
		PostalCode:          identityCard.PostalCode,
		Religion:            identityCard.Religion,
		Maritalstatus:       identityCard.Maritalstatus,
		Job:                 identityCard.Job,
		Citizenship:         identityCard.Citizenship,
		Validuntil:          identityCard.Validuntil,
		Cardphoto:           identityCard.Cardphoto,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}, nil
}

func (s *identityCardService) saveIdentityCardImage(userID string, cardPhoto *multipart.FileHeader) (string, error) {
	pathImage := "/uploads/identitycards/"
	cardPhotoDir := "./public" + os.Getenv("BASE_URL") + pathImage
	if _, err := os.Stat(cardPhotoDir); os.IsNotExist(err) {

		if err := os.MkdirAll(cardPhotoDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for identity card photo: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	extension := filepath.Ext(cardPhoto.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	cardPhotoFileName := fmt.Sprintf("%s_cardphoto%s", userID, extension)
	cardPhotoPath := filepath.Join(cardPhotoDir, cardPhotoFileName)

	src, err := cardPhoto.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(cardPhotoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create card photo file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save card photo: %v", err)
	}

	cardPhotoURL := fmt.Sprintf("%s%s", pathImage, cardPhotoFileName)

	return cardPhotoURL, nil
}

func deleteIdentityCardImage(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	baseDir := "./public/" + os.Getenv("BASE_URL")
	absolutePath := baseDir + imagePath

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %v", err)
	}

	err := os.Remove(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	log.Printf("Image deleted successfully: %s", absolutePath)
	return nil
}

func (s *identityCardService) CreateIdentityCard(ctx context.Context, userID string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*authentication.AuthResponse, error) {

	// Validate essential parameters
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	if request.DeviceID == "" {
		return nil, fmt.Errorf("deviceID cannot be empty")
	}

	cardPhotoPath, err := s.saveIdentityCardImage(userID, cardPhoto)
	if err != nil {
		return nil, fmt.Errorf("failed to save card photo: %v", err)
	}

	identityCard := &model.IdentityCard{
		UserID:              userID,
		Identificationumber: request.Identificationumber,
		Placeofbirth:        request.Placeofbirth,
		Dateofbirth:         request.Dateofbirth,
		Gender:              request.Gender,
		BloodType:           request.BloodType,
		Province:            request.Province,
		District:            request.District,
		SubDistrict:         request.SubDistrict,
		Hamlet:              request.Hamlet,
		Village:             request.Village,
		Neighbourhood:       request.Neighbourhood,
		PostalCode:          request.PostalCode,
		Religion:            request.Religion,
		Maritalstatus:       request.Maritalstatus,
		Job:                 request.Job,
		Citizenship:         request.Citizenship,
		Validuntil:          request.Validuntil,
		Cardphoto:           cardPhotoPath,
	}

	_, err = s.identityRepo.CreateIdentityCard(ctx, identityCard)
	if err != nil {
		log.Printf("Error creating identity card: %v", err)
		return nil, fmt.Errorf("failed to create identity card: %v", err)
	}

	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	// Validate user data
	if user.Role.RoleName == "" {
		return nil, fmt.Errorf("user role not found")
	}

	roleName := strings.ToLower(user.Role.RoleName)

	// Determine new registration status and progress
	var newRegistrationStatus string
	var newRegistrationProgress int

	switch roleName {
	case "pengepul":
		newRegistrationProgress = 2
		newRegistrationStatus = utils.RegStatusPending
	case "pengelola":
		newRegistrationProgress = 2
		newRegistrationStatus = user.RegistrationStatus
	default:
		newRegistrationProgress = int(user.RegistrationProgress)
		newRegistrationStatus = user.RegistrationStatus
	}

	// Update user registration progress and status
	updates := map[string]interface{}{
		"registration_progress": newRegistrationProgress,
		"registration_status":   newRegistrationStatus,
	}

	err = s.authRepo.PatchUser(ctx, userID, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Debug logging before token generation
	log.Printf("Token Generation Parameters:")
	log.Printf("- UserID: '%s'", user.ID)
	log.Printf("- Role: '%s'", user.Role.RoleName)
	log.Printf("- DeviceID: '%s'", request.DeviceID)
	log.Printf("- Registration Status: '%s'", newRegistrationStatus)

	// Generate token pair with updated status
	tokenResponse, err := utils.GenerateTokenPair(
		user.ID,
		user.Role.RoleName,
		request.DeviceID,
		newRegistrationStatus,
		newRegistrationProgress,
	)
	if err != nil {
		log.Printf("GenerateTokenPair error: %v", err)
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &authentication.AuthResponse{
		Message:            "identity card berhasil diunggah, silakan tunggu konfirmasi dari admin dalam 1x24 jam",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: newRegistrationStatus,
		NextStep:           tokenResponse.NextStep,
		SessionID:          tokenResponse.SessionID,
	}, nil
}

func (s *identityCardService) GetIdentityCardByID(ctx context.Context, id string) (*ResponseIdentityCardDTO, error) {
	identityCard, err := s.identityRepo.GetIdentityCardByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching identity card: %v", err)
		return nil, fmt.Errorf("failed to fetch identity card")
	}
	return FormatResponseIdentityCard(identityCard)
}

func (s *identityCardService) GetIdentityCardsByUserID(ctx context.Context, userID string) ([]ResponseIdentityCardDTO, error) {
	identityCards, err := s.identityRepo.GetIdentityCardsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error fetching identity cards by userID: %v", err)
		return nil, fmt.Errorf("failed to fetch identity cards by userID")
	}

	var response []ResponseIdentityCardDTO
	for _, card := range identityCards {
		dto, _ := FormatResponseIdentityCard(&card)
		response = append(response, *dto)
	}
	return response, nil
}

func (s *identityCardService) UpdateIdentityCard(ctx context.Context, userID string, id string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*ResponseIdentityCardDTO, error) {

	errors, valid := request.ValidateIdentityCardInput()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	identityCard, err := s.identityRepo.GetIdentityCardByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("identity card not found: %v", err)
	}

	if identityCard.Cardphoto != "" {
		err := deleteIdentityCardImage(identityCard.Cardphoto)
		if err != nil {
			return nil, fmt.Errorf("failed to delete old image: %v", err)
		}
	}

	var cardPhotoPath string
	if cardPhoto != nil {
		cardPhotoPath, err = s.saveIdentityCardImage(userID, cardPhoto)
		if err != nil {
			return nil, fmt.Errorf("failed to save card photo: %v", err)
		}
	}

	identityCard.Identificationumber = request.Identificationumber
	identityCard.Placeofbirth = request.Placeofbirth
	identityCard.Dateofbirth = request.Dateofbirth
	identityCard.Gender = request.Gender
	identityCard.BloodType = request.BloodType
	identityCard.Province = request.Province
	identityCard.District = request.District
	identityCard.SubDistrict = request.SubDistrict
	identityCard.Hamlet = request.Hamlet
	identityCard.Village = request.Village
	identityCard.Neighbourhood = request.Neighbourhood
	identityCard.PostalCode = request.PostalCode
	identityCard.Religion = request.Religion
	identityCard.Maritalstatus = request.Maritalstatus
	identityCard.Job = request.Job
	identityCard.Citizenship = request.Citizenship
	identityCard.Validuntil = request.Validuntil
	if cardPhotoPath != "" {
		identityCard.Cardphoto = cardPhotoPath
	}

	if err != nil {
		log.Printf("Error updating identity card: %v", err)
		return nil, fmt.Errorf("failed to update identity card: %v", err)
	}

	idcardResponseDTO, _ := FormatResponseIdentityCard(identityCard)

	return idcardResponseDTO, nil
}
