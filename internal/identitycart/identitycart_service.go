package identitycart

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rijig/internal/authentication"
	"rijig/internal/role"
	"rijig/internal/userprofile"
	"rijig/model"
	"rijig/utils"
	"time"
)

type IdentityCardService interface {
	CreateIdentityCard(ctx context.Context, userID, deviceID string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*authentication.AuthResponse, error)
	GetIdentityCardByID(ctx context.Context, id string) (*ResponseIdentityCardDTO, error)
	GetIdentityCardsByUserID(ctx context.Context, userID string) ([]ResponseIdentityCardDTO, error)
	UpdateIdentityCard(ctx context.Context, userID string, id string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*ResponseIdentityCardDTO, error)

	GetAllIdentityCardsByRegStatus(ctx context.Context, userRegStatus string) ([]ResponseIdentityCardDTO, error)
	UpdateUserRegistrationStatusByIdentityCard(ctx context.Context, identityCardUserID string, newStatus string) error
}

type identityCardService struct {
	identityRepo IdentityCardRepository
	authRepo     authentication.AuthenticationRepository
	userRepo     userprofile.UserProfileRepository
}

func NewIdentityCardService(identityRepo IdentityCardRepository, authRepo authentication.AuthenticationRepository, userRepo userprofile.UserProfileRepository) IdentityCardService {
	return &identityCardService{
		identityRepo,
		authRepo, userRepo,
	}
}

type IdentityCardWithUserDTO struct {
	IdentityCard ResponseIdentityCardDTO            `json:"identity_card"`
	User         userprofile.UserProfileResponseDTO `json:"user"`
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

	if _, err := io.Copy(dst, src); err != nil {
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

func (s *identityCardService) CreateIdentityCard(ctx context.Context, userID, deviceID string, request *RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*authentication.AuthResponse, error) {

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

	if user.Role.RoleName == "" {
		return nil, fmt.Errorf("user role not found")
	}

	updates := map[string]interface{}{
		"registration_progress": utils.ProgressDataSubmitted,
		"registration_status":   utils.RegStatusPending,
	}

	err = s.authRepo.PatchUser(ctx, userID, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	updated, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userprofile.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	log.Printf("Token Generation Parameters:")
	log.Printf("- UserID: '%s'", user.ID)
	log.Printf("- Role: '%s'", user.Role.RoleName)
	log.Printf("- DeviceID: '%s'", deviceID)
	log.Printf("- Registration Status: '%s'", utils.RegStatusPending)

	tokenResponse, err := utils.GenerateTokenPair(
		updated.ID,
		updated.Role.RoleName,
		deviceID,
		updated.RegistrationStatus,
		int(updated.RegistrationProgress),
	)
	if err != nil {
		log.Printf("GenerateTokenPair error: %v", err)
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	nextStep := utils.GetNextRegistrationStep(
		updated.Role.RoleName,
		int(updated.RegistrationProgress),
		updated.RegistrationStatus,
	)

	return &authentication.AuthResponse{
		Message:            "identity card berhasil diunggah, silakan tunggu konfirmasi dari admin dalam 1x24 jam",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: updated.RegistrationStatus,
		NextStep:           nextStep,
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

func (s *identityCardService) GetAllIdentityCardsByRegStatus(ctx context.Context, userRegStatus string) ([]ResponseIdentityCardDTO, error) {
	identityCards, err := s.authRepo.GetIdentityCardsByUserRegStatus(ctx, userRegStatus)
	if err != nil {
		log.Printf("Error getting identity cards by registration status: %v", err)
		return nil, fmt.Errorf("failed to get identity cards: %w", err)
	}

	var response []ResponseIdentityCardDTO
	for _, card := range identityCards {
		createdAt, _ := utils.FormatDateToIndonesianFormat(card.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(card.UpdatedAt)
		dto := ResponseIdentityCardDTO{
			ID:                  card.ID,
			UserID:              card.UserID,
			Identificationumber: card.Identificationumber,
			Placeofbirth:        card.Placeofbirth,
			Dateofbirth:         card.Dateofbirth,
			Gender:              card.Gender,
			BloodType:           card.BloodType,
			Province:            card.Province,
			District:            card.District,
			SubDistrict:         card.SubDistrict,
			Hamlet:              card.Hamlet,
			Village:             card.Village,
			Neighbourhood:       card.Neighbourhood,
			PostalCode:          card.PostalCode,
			Religion:            card.Religion,
			Maritalstatus:       card.Maritalstatus,
			Job:                 card.Job,
			Citizenship:         card.Citizenship,
			Validuntil:          card.Validuntil,
			Cardphoto:           card.Cardphoto,
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
		}
		response = append(response, dto)
	}

	return response, nil
}

func (s *identityCardService) UpdateUserRegistrationStatusByIdentityCard(ctx context.Context, identityCardUserID string, newStatus string) error {

	user, err := s.authRepo.FindUserByID(ctx, identityCardUserID)
	if err != nil {
		log.Printf("Error finding user by ID %s: %v", identityCardUserID, err)
		return fmt.Errorf("user not found: %w", err)
	}

	updates := map[string]interface{}{
		"registration_status": newStatus,
		"updated_at":          time.Now(),
	}

	switch newStatus {
	case utils.RegStatusConfirmed:
		updates["registration_progress"] = utils.ProgressDataSubmitted

		identityCards, err := s.GetIdentityCardsByUserID(ctx, identityCardUserID)
		if err != nil {
			log.Printf("Error fetching identity cards for user ID %s: %v", identityCardUserID, err)
			return fmt.Errorf("failed to fetch identity card data: %w", err)
		}

		if len(identityCards) == 0 {
			log.Printf("No identity card found for user ID %s", identityCardUserID)
			return fmt.Errorf("no identity card found for user")
		}

		identityCard := identityCards[0]

		updates["name"] = identityCard.Fullname
		updates["gender"] = identityCard.Gender
		updates["dateofbirth"] = identityCard.Dateofbirth
		updates["placeofbirth"] = identityCard.District

		log.Printf("Syncing user data for ID %s: name=%s, gender=%s, dob=%s, pob=%s",
			identityCardUserID, identityCard.Fullname, identityCard.Gender,
			identityCard.Dateofbirth, identityCard.District)

	case utils.RegStatusRejected:
		updates["registration_progress"] = utils.ProgressOTPVerified

	}

	err = s.authRepo.PatchUser(ctx, user.ID, updates)
	if err != nil {
		log.Printf("Error updating user registration status for user ID %s: %v", user.ID, err)
		return fmt.Errorf("failed to update user registration status: %w", err)
	}

	log.Printf("Successfully updated registration status for user ID %s to %s", user.ID, newStatus)

	if newStatus == utils.RegStatusConfirmed {
		log.Printf("User profile data synced successfully for user ID %s", user.ID)
	}

	return nil
}

func (s *identityCardService) mapIdentityCardToDTO(card model.IdentityCard) ResponseIdentityCardDTO {
	createdAt, _ := utils.FormatDateToIndonesianFormat(card.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(card.UpdatedAt)
	return ResponseIdentityCardDTO{
		ID:                  card.ID,
		UserID:              card.UserID,
		Identificationumber: card.Identificationumber,
		Placeofbirth:        card.Placeofbirth,
		Dateofbirth:         card.Dateofbirth,
		Gender:              card.Gender,
		BloodType:           card.BloodType,
		Province:            card.Province,
		District:            card.District,
		SubDistrict:         card.SubDistrict,
		Hamlet:              card.Hamlet,
		Village:             card.Village,
		Neighbourhood:       card.Neighbourhood,
		PostalCode:          card.PostalCode,
		Religion:            card.Religion,
		Maritalstatus:       card.Maritalstatus,
		Job:                 card.Job,
		Citizenship:         card.Citizenship,
		Validuntil:          card.Validuntil,
		Cardphoto:           card.Cardphoto,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

func (s *identityCardService) mapUserToDTO(user model.User) userprofile.UserProfileResponseDTO {
	avatar := ""
	if user.Avatar != nil {
		avatar = *user.Avatar
	}

	var roleDTO role.RoleResponseDTO
	if user.Role != nil {
		roleDTO = role.RoleResponseDTO{
			ID:       user.Role.ID,
			RoleName: user.Role.RoleName,
		}
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)
	return userprofile.UserProfileResponseDTO{
		ID:            user.ID,
		Avatar:        avatar,
		Name:          user.Name,
		Gender:        user.Gender,
		Dateofbirth:   user.Dateofbirth,
		Placeofbirth:  user.Placeofbirth,
		Phone:         user.Phone,
		Email:         user.Email,
		PhoneVerified: user.PhoneVerified,
		Role:          roleDTO,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
