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

	"github.com/google/uuid"
)

type AboutService interface {
	CreateAbout(request dto.RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*dto.ResponseAboutDTO, error)
	UpdateAbout(id string, request dto.RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*dto.ResponseAboutDTO, error)
	GetAllAbout() ([]dto.ResponseAboutDTO, error)
	GetAboutByID(id string) (*dto.ResponseAboutDTO, error)
	GetAboutDetailById(id string) (*dto.ResponseAboutDetailDTO, error)
	DeleteAbout(id string) error

	CreateAboutDetail(request dto.RequestAboutDetailDTO, coverImageAboutDetail *multipart.FileHeader) (*dto.ResponseAboutDetailDTO, error)
	UpdateAboutDetail(id string, request dto.RequestAboutDetailDTO, imageDetail *multipart.FileHeader) (*dto.ResponseAboutDetailDTO, error)
	DeleteAboutDetail(id string) error
}

type aboutService struct {
	aboutRepo repositories.AboutRepository
}

func NewAboutService(aboutRepo repositories.AboutRepository) AboutService {
	return &aboutService{aboutRepo: aboutRepo}
}

func formatResponseAboutDetailDTO(about *model.AboutDetail) (*dto.ResponseAboutDetailDTO, error) {
	createdAt, _ := utils.FormatDateToIndonesianFormat(about.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(about.UpdatedAt)

	response := &dto.ResponseAboutDetailDTO{
		ID:          about.ID,
		AboutID:     about.AboutID,
		ImageDetail: about.ImageDetail,
		Description: about.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return response, nil
}

func formatResponseAboutDTO(about *model.About) (*dto.ResponseAboutDTO, error) {
	createdAt, _ := utils.FormatDateToIndonesianFormat(about.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(about.UpdatedAt)

	response := &dto.ResponseAboutDTO{
		ID:         about.ID,
		Title:      about.Title,
		CoverImage: about.CoverImage,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	return response, nil
}

func (s *aboutService) saveCoverImageAbout(coverImageAbout *multipart.FileHeader) (string, error) {
	pathImage := "/uploads/coverabout/"
	coverImageAboutDir := "./public" + os.Getenv("BASE_URL") + pathImage
	if _, err := os.Stat(coverImageAboutDir); os.IsNotExist(err) {

		if err := os.MkdirAll(coverImageAboutDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("gagal membuat direktori untuk cover image about: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := filepath.Ext(coverImageAbout.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	coverImageFileName := fmt.Sprintf("%s_coverabout%s", uuid.New().String(), extension)
	coverImagePath := filepath.Join(coverImageAboutDir, coverImageFileName)

	src, err := coverImageAbout.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(coverImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cover image about file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save cover image about: %v", err)
	}

	coverImageAboutUrl := fmt.Sprintf("%s%s", pathImage, coverImageFileName)

	return coverImageAboutUrl, nil
}

func (s *aboutService) saveCoverImageAboutDetail(coverImageAbout *multipart.FileHeader) (string, error) {
	pathImage := "/uploads/coverabout/coveraboutdetail/"
	coverImageAboutDir := "./public" + os.Getenv("BASE_URL") + pathImage
	if _, err := os.Stat(coverImageAboutDir); os.IsNotExist(err) {

		if err := os.MkdirAll(coverImageAboutDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("gagal membuat direktori untuk cover image about: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := filepath.Ext(coverImageAbout.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	coverImageFileName := fmt.Sprintf("%s_coveraboutdetail_%s", uuid.New().String(), extension)
	coverImagePath := filepath.Join(coverImageAboutDir, coverImageFileName)

	src, err := coverImageAbout.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(coverImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cover image about file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save cover image about: %v", err)
	}

	coverImageAboutUrl := fmt.Sprintf("%s%s", pathImage, coverImageFileName)

	return coverImageAboutUrl, nil
}

func deleteCoverImageAbout(coverimageAboutPath string) error {
	if coverimageAboutPath == "" {
		return nil
	}

	baseDir := "./public/" + os.Getenv("BASE_URL")
	absolutePath := baseDir + coverimageAboutPath

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

func (s *aboutService) CreateAbout(request dto.RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*dto.ResponseAboutDTO, error) {
	errors, valid := request.ValidateAbout()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	coverImageAboutPath, err := s.saveCoverImageAbout(coverImageAbout)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan cover image about: %v ", err)
	}

	about := model.About{
		Title:      request.Title,
		CoverImage: coverImageAboutPath,
	}

	if err := s.aboutRepo.CreateAbout(&about); err != nil {
		return nil, fmt.Errorf("failed to create About: %v", err)
	}

	response, err := formatResponseAboutDTO(&about)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	return response, nil
}

func (s *aboutService) UpdateAbout(id string, request dto.RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*dto.ResponseAboutDTO, error) {

	errors, valid := request.ValidateAbout()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	about, err := s.aboutRepo.GetAboutByID(id)
	if err != nil {
		return nil, fmt.Errorf("about not found: %v", err)
	}

	if about.CoverImage != "" {
		err := deleteCoverImageAbout(about.CoverImage)
		if err != nil {
			return nil, fmt.Errorf("gagal mengahpus gambar lama: %v", err)
		}
	}

	var coverImageAboutPath string
	if coverImageAbout != nil {
		coverImageAboutPath, err = s.saveCoverImageAbout(coverImageAbout)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan gambar baru: %v", err)
		}
	}

	about.Title = request.Title
	if coverImageAboutPath != "" {
		about.CoverImage = coverImageAboutPath
	}

	updatedAbout, err := s.aboutRepo.UpdateAbout(id, about)
	if err != nil {
		return nil, fmt.Errorf("failed to update About: %v", err)
	}

	response, err := formatResponseAboutDTO(updatedAbout)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	return response, nil
}

func (s *aboutService) GetAllAbout() ([]dto.ResponseAboutDTO, error) {

	aboutList, err := s.aboutRepo.GetAllAbout()
	if err != nil {
		return nil, fmt.Errorf("failed to get About list: %v", err)
	}

	var aboutDTOList []dto.ResponseAboutDTO
	for _, about := range aboutList {
		response, err := formatResponseAboutDTO(&about)
		if err != nil {
			log.Printf("Error formatting About response: %v", err)
			continue
		}
		aboutDTOList = append(aboutDTOList, *response)
	}

	return aboutDTOList, nil
}

func (s *aboutService) GetAboutByID(id string) (*dto.ResponseAboutDTO, error) {

	about, err := s.aboutRepo.GetAboutByID(id)
	if err != nil {
		return nil, fmt.Errorf("about not found: %v", err)
	}

	response, err := formatResponseAboutDTO(about)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	var responseDetails []dto.ResponseAboutDetailDTO
	for _, detail := range about.AboutDetail {
		formattedDetail, err := formatResponseAboutDetailDTO(&detail)
		if err != nil {
			return nil, fmt.Errorf("error formatting AboutDetail response: %v", err)
		}
		responseDetails = append(responseDetails, *formattedDetail)
	}

	response.AboutDetail = &responseDetails

	return response, nil
}

func (s *aboutService) GetAboutDetailById(id string) (*dto.ResponseAboutDetailDTO, error) {

	about, err := s.aboutRepo.GetAboutDetailByID(id)
	if err != nil {
		return nil, fmt.Errorf("about not found: %v", err)
	}

	response, err := formatResponseAboutDetailDTO(about)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	return response, nil
}

func (s *aboutService) DeleteAbout(id string) error {
	about, err := s.aboutRepo.GetAboutByID(id)
	if err != nil {
		return fmt.Errorf("about not found: %v", err)
	}

	if about.CoverImage != "" {
		err := deleteCoverImageAbout(about.CoverImage)
		if err != nil {
			return fmt.Errorf("gagal mengahpus gambar lama: %v", err)
		}
	}

	if err := s.aboutRepo.DeleteAbout(id); err != nil {
		return fmt.Errorf("failed to delete About: %v", err)
	}

	return nil
}

func (s *aboutService) CreateAboutDetail(request dto.RequestAboutDetailDTO, coverImageAboutDetail *multipart.FileHeader) (*dto.ResponseAboutDetailDTO, error) {

	errors, valid := request.ValidateAboutDetail()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	_, err := s.aboutRepo.GetAboutByIDWithoutPrel(request.AboutId)
	if err != nil {
		return nil, fmt.Errorf("about_id tidak ditemukan: %v", err)
	}

	coverImageAboutDetailPath, err := s.saveCoverImageAboutDetail(coverImageAboutDetail)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan cover image about detail: %v ", err)
	}

	aboutDetail := model.AboutDetail{
		AboutID:     request.AboutId,
		ImageDetail: coverImageAboutDetailPath,
		Description: request.Description,
	}

	if err := s.aboutRepo.CreateAboutDetail(&aboutDetail); err != nil {
		return nil, fmt.Errorf("failed to create AboutDetail: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(aboutDetail.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(aboutDetail.UpdatedAt)

	response := &dto.ResponseAboutDetailDTO{
		ID:          aboutDetail.ID,
		AboutID:     aboutDetail.AboutID,
		ImageDetail: aboutDetail.ImageDetail,
		Description: aboutDetail.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return response, nil
}

func (s *aboutService) UpdateAboutDetail(id string, request dto.RequestAboutDetailDTO, imageDetail *multipart.FileHeader) (*dto.ResponseAboutDetailDTO, error) {

	errors, valid := request.ValidateAboutDetail()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	aboutDetail, err := s.aboutRepo.GetAboutDetailByID(id)
	if err != nil {
		return nil, fmt.Errorf("about detail tidakck ditemukan: %v", err)
	}

	if aboutDetail.ImageDetail != "" {
		err := deleteCoverImageAbout(aboutDetail.ImageDetail)
		if err != nil {
			return nil, fmt.Errorf("gagal menghapus gambar lama: %v", err)
		}
	}

	var coverImageAboutDeatilPath string
	if imageDetail != nil {
		coverImageAboutDeatilPath, err = s.saveCoverImageAbout(imageDetail)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan gambar baru: %v", err)
		}
	}

	aboutDetail.Description = request.Description
	if coverImageAboutDeatilPath != "" {
		aboutDetail.ImageDetail = coverImageAboutDeatilPath
	}

	aboutDetail, err = s.aboutRepo.UpdateAboutDetail(id, aboutDetail)
	if err != nil {
		log.Printf("Error updating about detail: %v", err)
		return nil, fmt.Errorf("failed to update about detail: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(aboutDetail.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(aboutDetail.UpdatedAt)

	response := &dto.ResponseAboutDetailDTO{
		ID:          aboutDetail.ID,
		AboutID:     aboutDetail.AboutID,
		ImageDetail: aboutDetail.ImageDetail,
		Description: aboutDetail.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return response, nil
}

func (s *aboutService) DeleteAboutDetail(id string) error {
	aboutDetail, err := s.aboutRepo.GetAboutDetailByID(id)
	if err != nil {
		return fmt.Errorf("about detail tidakck ditemukan: %v", err)
	}

	if aboutDetail.ImageDetail != "" {
		err := deleteCoverImageAbout(aboutDetail.ImageDetail)
		if err != nil {
			return fmt.Errorf("gagal menghapus gambar lama: %v", err)
		}
	}

	if err := s.aboutRepo.DeleteAboutDetail(id); err != nil {
		return fmt.Errorf("failed to delete AboutDetail: %v", err)
	}
	return nil
}
