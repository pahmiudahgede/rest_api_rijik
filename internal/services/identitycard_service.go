package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
)

type IdentityCardService interface {
	CreateIdentityCard(userID string, request *dto.RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*dto.ResponseIdentityCardDTO, error)
	GetIdentityCardByID(id string) (*dto.ResponseIdentityCardDTO, error)
	GetIdentityCardsByUserID(userID string) ([]dto.ResponseIdentityCardDTO, error)
	UpdateIdentityCard(userID string, id string, request *dto.RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*dto.ResponseIdentityCardDTO, error)
	DeleteIdentityCard(id string) error
}

type identityCardService struct {
	identityCardRepo repositories.IdentityCardRepository
	userRepo         repositories.UserProfilRepository
}

func NewIdentityCardService(identityCardRepo repositories.IdentityCardRepository, userRepo repositories.UserProfilRepository) IdentityCardService {
	return &identityCardService{
		identityCardRepo: identityCardRepo,
		userRepo:         userRepo,
	}
}

func FormatResponseIdentityCars(identityCard *model.IdentityCard) (*dto.ResponseIdentityCardDTO, error) {

	createdAt, _ := utils.FormatDateToIndonesianFormat(identityCard.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(identityCard.UpdatedAt)

	idcardResponseDTO := &dto.ResponseIdentityCardDTO{
		ID:                  identityCard.ID,
		UserID:              identityCard.UserID,
		Identificationumber: identityCard.Identificationumber,
		Placeofbirth:        identityCard.Placeofbirth,
		Dateofbirth:         identityCard.Dateofbirth,
		Gender:              identityCard.Gender,
		BloodType:           identityCard.BloodType,
		District:            identityCard.District,
		Village:             identityCard.Village,
		Neighbourhood:       identityCard.Neighbourhood,
		Religion:            identityCard.Religion,
		Maritalstatus:       identityCard.Maritalstatus,
		Job:                 identityCard.Job,
		Citizenship:         identityCard.Citizenship,
		Validuntil:          identityCard.Validuntil,
		Cardphoto:           identityCard.Cardphoto,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}

	return idcardResponseDTO, nil
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

func (s *identityCardService) CreateIdentityCard(userID string, request *dto.RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*dto.ResponseIdentityCardDTO, error) {

	errors, valid := request.ValidateIdentityCardInput()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
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
		District:            request.District,
		Village:             request.Village,
		Neighbourhood:       request.Neighbourhood,
		Religion:            request.Religion,
		Maritalstatus:       request.Maritalstatus,
		Job:                 request.Job,
		Citizenship:         request.Citizenship,
		Validuntil:          request.Validuntil,
		Cardphoto:           cardPhotoPath,
	}

	identityCard, err = s.identityCardRepo.CreateIdentityCard(identityCard)
	if err != nil {
		log.Printf("Error creating identity card: %v", err)
		return nil, fmt.Errorf("failed to create identity card: %v", err)
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failde to fint user: %v", err)
	}

	user.RegistrationStatus = "onreview"

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}
	idcardResponseDTO, _ := FormatResponseIdentityCars(identityCard)

	return idcardResponseDTO, nil
}

func (s *identityCardService) GetIdentityCardByID(id string) (*dto.ResponseIdentityCardDTO, error) {

	identityCard, err := s.identityCardRepo.GetIdentityCardByID(id)
	if err != nil {
		log.Printf("Error fetching identity card: %v", err)
		return nil, fmt.Errorf("failed to fetch identity card")
	}

	idcardResponseDTO, _ := FormatResponseIdentityCars(identityCard)

	return idcardResponseDTO, nil

}

func (s *identityCardService) GetIdentityCardsByUserID(userID string) ([]dto.ResponseIdentityCardDTO, error) {

	identityCards, err := s.identityCardRepo.GetIdentityCardsByUserID(userID)
	if err != nil {
		log.Printf("Error fetching identity cards by userID: %v", err)
		return nil, fmt.Errorf("failed to fetch identity cards by userID")
	}

	var response []dto.ResponseIdentityCardDTO
	for _, card := range identityCards {

		idcardResponseDTO, err := FormatResponseIdentityCars(&card)
		if err != nil {
			log.Printf("Error creating response DTO for identity card ID %v: %v", card.ID, err)

			continue
		}
		response = append(response, *idcardResponseDTO)
	}

	return response, nil
}

func (s *identityCardService) UpdateIdentityCard(userID string, id string, request *dto.RequestIdentityCardDTO, cardPhoto *multipart.FileHeader) (*dto.ResponseIdentityCardDTO, error) {

	errors, valid := request.ValidateIdentityCardInput()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	identityCard, err := s.identityCardRepo.GetIdentityCardByID(id)
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
	identityCard.District = request.District
	identityCard.Village = request.Village
	identityCard.Neighbourhood = request.Neighbourhood
	identityCard.Religion = request.Religion
	identityCard.Maritalstatus = request.Maritalstatus
	identityCard.Job = request.Job
	identityCard.Citizenship = request.Citizenship
	identityCard.Validuntil = request.Validuntil
	if cardPhotoPath != "" {
		identityCard.Cardphoto = cardPhotoPath
	}

	identityCard, err = s.identityCardRepo.UpdateIdentityCard(id, identityCard)
	if err != nil {
		log.Printf("Error updating identity card: %v", err)
		return nil, fmt.Errorf("failed to update identity card: %v", err)
	}

	idcardResponseDTO, _ := FormatResponseIdentityCars(identityCard)

	return idcardResponseDTO, nil
}

func (s *identityCardService) DeleteIdentityCard(id string) error {

	identityCard, err := s.identityCardRepo.GetIdentityCardByID(id)
	if err != nil {
		return fmt.Errorf("identity card not found: %v", err)
	}

	if identityCard.Cardphoto != "" {
		err := deleteIdentityCardImage(identityCard.Cardphoto)
		if err != nil {
			return fmt.Errorf("failed to delete card photo: %v", err)
		}
	}

	err = s.identityCardRepo.DeleteIdentityCard(id)
	if err != nil {
		return fmt.Errorf("failed to delete identity card: %v", err)
	}

	return nil
}
