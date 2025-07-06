package about

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rijig/model"
	"rijig/utils"
	"time"

	"github.com/google/uuid"
)

const (
	cacheKeyAllAbout    = "about:all"
	cacheKeyAboutByID   = "about:id:%s"
	cacheKeyAboutDetail = "about_detail:id:%s"

	cacheTTL = 30 * time.Minute
)

type AboutService interface {
	CreateAbout(ctx context.Context, request RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*ResponseAboutDTO, error)
	UpdateAbout(ctx context.Context, id string, request RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*ResponseAboutDTO, error)
	GetAllAbout(ctx context.Context) ([]ResponseAboutDTO, error)
	GetAboutByID(ctx context.Context, id string) (*ResponseAboutDTO, error)
	GetAboutDetailById(ctx context.Context, id string) (*ResponseAboutDetailDTO, error)
	DeleteAbout(ctx context.Context, id string) error

	CreateAboutDetail(ctx context.Context, request RequestAboutDetailDTO, coverImageAboutDetail *multipart.FileHeader) (*ResponseAboutDetailDTO, error)
	UpdateAboutDetail(ctx context.Context, id string, request RequestAboutDetailDTO, imageDetail *multipart.FileHeader) (*ResponseAboutDetailDTO, error)
	DeleteAboutDetail(ctx context.Context, id string) error
}

type aboutService struct {
	aboutRepo AboutRepository
}

func NewAboutService(aboutRepo AboutRepository) AboutService {
	return &aboutService{aboutRepo: aboutRepo}
}

func (s *aboutService) invalidateAboutCaches(aboutID string) {

	if err := utils.DeleteCache(cacheKeyAllAbout); err != nil {
		log.Printf("Failed to invalidate all about cache: %v", err)
	}

	aboutCacheKey := fmt.Sprintf(cacheKeyAboutByID, aboutID)
	if err := utils.DeleteCache(aboutCacheKey); err != nil {
		log.Printf("Failed to invalidate about cache for ID %s: %v", aboutID, err)
	}
}

func (s *aboutService) invalidateAboutDetailCaches(aboutDetailID, aboutID string) {

	detailCacheKey := fmt.Sprintf(cacheKeyAboutDetail, aboutDetailID)
	if err := utils.DeleteCache(detailCacheKey); err != nil {
		log.Printf("Failed to invalidate about detail cache for ID %s: %v", aboutDetailID, err)
	}

	s.invalidateAboutCaches(aboutID)
}

func formatResponseAboutDetailDTO(about *model.AboutDetail) (*ResponseAboutDetailDTO, error) {
	createdAt, _ := utils.FormatDateToIndonesianFormat(about.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(about.UpdatedAt)

	response := &ResponseAboutDetailDTO{
		ID:          about.ID,
		AboutID:     about.AboutID,
		ImageDetail: about.ImageDetail,
		Description: about.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return response, nil
}

func formatResponseAboutDTO(about *model.About) (*ResponseAboutDTO, error) {
	createdAt, _ := utils.FormatDateToIndonesianFormat(about.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(about.UpdatedAt)

	response := &ResponseAboutDTO{
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
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, .png, and .svg are allowed")
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
			return "", fmt.Errorf("gagal membuat direktori untuk cover image about detail: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := filepath.Ext(coverImageAbout.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, .png, and .svg are allowed")
	}

	coverImageFileName := fmt.Sprintf("%s_coveraboutdetail%s", uuid.New().String(), extension)
	coverImagePath := filepath.Join(coverImageAboutDir, coverImageFileName)

	src, err := coverImageAbout.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(coverImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cover image about detail file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save cover image about detail: %v", err)
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

func (s *aboutService) CreateAbout(ctx context.Context, request RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*ResponseAboutDTO, error) {
	errors, valid := request.ValidateAbout()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	coverImageAboutPath, err := s.saveCoverImageAbout(coverImageAbout)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan cover image about: %v", err)
	}

	about := model.About{
		Title:      request.Title,
		CoverImage: coverImageAboutPath,
	}

	if err := s.aboutRepo.CreateAbout(ctx, &about); err != nil {
		return nil, fmt.Errorf("failed to create About: %v", err)
	}

	response, err := formatResponseAboutDTO(&about)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	s.invalidateAboutCaches("")

	return response, nil
}

func (s *aboutService) UpdateAbout(ctx context.Context, id string, request RequestAboutDTO, coverImageAbout *multipart.FileHeader) (*ResponseAboutDTO, error) {
	errors, valid := request.ValidateAbout()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	about, err := s.aboutRepo.GetAboutByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("about not found: %v", err)
	}

	oldCoverImage := about.CoverImage

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

	updatedAbout, err := s.aboutRepo.UpdateAbout(ctx, id, about)
	if err != nil {
		return nil, fmt.Errorf("failed to update About: %v", err)
	}

	if oldCoverImage != "" && coverImageAboutPath != "" {
		if err := deleteCoverImageAbout(oldCoverImage); err != nil {
			log.Printf("Warning: failed to delete old image: %v", err)
		}
	}

	response, err := formatResponseAboutDTO(updatedAbout)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	s.invalidateAboutCaches(id)

	return response, nil
}

func (s *aboutService) GetAllAbout(ctx context.Context) ([]ResponseAboutDTO, error) {

	var cachedAbouts []ResponseAboutDTO
	if err := utils.GetCache(cacheKeyAllAbout, &cachedAbouts); err == nil {
		return cachedAbouts, nil
	}

	aboutList, err := s.aboutRepo.GetAllAbout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get About list: %v", err)
	}

	var aboutDTOList []ResponseAboutDTO
	for _, about := range aboutList {
		response, err := formatResponseAboutDTO(&about)
		if err != nil {
			log.Printf("Error formatting About response: %v", err)
			continue
		}
		aboutDTOList = append(aboutDTOList, *response)
	}

	if err := utils.SetCache(cacheKeyAllAbout, aboutDTOList, cacheTTL); err != nil {
		log.Printf("Failed to cache all about data: %v", err)
	}

	return aboutDTOList, nil
}

func (s *aboutService) GetAboutByID(ctx context.Context, id string) (*ResponseAboutDTO, error) {
	cacheKey := fmt.Sprintf(cacheKeyAboutByID, id)

	var cachedAbout ResponseAboutDTO
	if err := utils.GetCache(cacheKey, &cachedAbout); err == nil {
		return &cachedAbout, nil
	}

	about, err := s.aboutRepo.GetAboutByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("about not found: %v", err)
	}

	response, err := formatResponseAboutDTO(about)
	if err != nil {
		return nil, fmt.Errorf("error formatting About response: %v", err)
	}

	var responseDetails []ResponseAboutDetailDTO
	for _, detail := range about.AboutDetail {
		formattedDetail, err := formatResponseAboutDetailDTO(&detail)
		if err != nil {
			return nil, fmt.Errorf("error formatting AboutDetail response: %v", err)
		}
		responseDetails = append(responseDetails, *formattedDetail)
	}

	response.AboutDetail = &responseDetails

	if err := utils.SetCache(cacheKey, response, cacheTTL); err != nil {
		log.Printf("Failed to cache about data for ID %s: %v", id, err)
	}

	return response, nil
}

func (s *aboutService) GetAboutDetailById(ctx context.Context, id string) (*ResponseAboutDetailDTO, error) {
	cacheKey := fmt.Sprintf(cacheKeyAboutDetail, id)

	var cachedDetail ResponseAboutDetailDTO
	if err := utils.GetCache(cacheKey, &cachedDetail); err == nil {
		return &cachedDetail, nil
	}

	aboutDetail, err := s.aboutRepo.GetAboutDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("about detail not found: %v", err)
	}

	response, err := formatResponseAboutDetailDTO(aboutDetail)
	if err != nil {
		return nil, fmt.Errorf("error formatting AboutDetail response: %v", err)
	}

	if err := utils.SetCache(cacheKey, response, cacheTTL); err != nil {
		log.Printf("Failed to cache about detail data for ID %s: %v", id, err)
	}

	return response, nil
}

func (s *aboutService) DeleteAbout(ctx context.Context, id string) error {
	about, err := s.aboutRepo.GetAboutByID(ctx, id)
	if err != nil {
		return fmt.Errorf("about not found: %v", err)
	}

	if about.CoverImage != "" {
		if err := deleteCoverImageAbout(about.CoverImage); err != nil {
			log.Printf("Warning: failed to delete cover image: %v", err)
		}
	}

	for _, detail := range about.AboutDetail {
		if detail.ImageDetail != "" {
			if err := deleteCoverImageAbout(detail.ImageDetail); err != nil {
				log.Printf("Warning: failed to delete detail image: %v", err)
			}
		}
	}

	if err := s.aboutRepo.DeleteAbout(ctx, id); err != nil {
		return fmt.Errorf("failed to delete About: %v", err)
	}

	s.invalidateAboutCaches(id)

	return nil
}

func (s *aboutService) CreateAboutDetail(ctx context.Context, request RequestAboutDetailDTO, coverImageAboutDetail *multipart.FileHeader) (*ResponseAboutDetailDTO, error) {
	errors, valid := request.ValidateAboutDetail()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	_, err := s.aboutRepo.GetAboutByIDWithoutPrel(ctx, request.AboutId)
	if err != nil {
		return nil, fmt.Errorf("about_id tidak ditemukan: %v", err)
	}

	coverImageAboutDetailPath, err := s.saveCoverImageAboutDetail(coverImageAboutDetail)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan cover image about detail: %v", err)
	}

	aboutDetail := model.AboutDetail{
		AboutID:     request.AboutId,
		ImageDetail: coverImageAboutDetailPath,
		Description: request.Description,
	}

	if err := s.aboutRepo.CreateAboutDetail(ctx, &aboutDetail); err != nil {
		return nil, fmt.Errorf("failed to create AboutDetail: %v", err)
	}

	response, err := formatResponseAboutDetailDTO(&aboutDetail)
	if err != nil {
		return nil, fmt.Errorf("error formatting AboutDetail response: %v", err)
	}

	s.invalidateAboutDetailCaches("", request.AboutId)

	return response, nil
}

func (s *aboutService) UpdateAboutDetail(ctx context.Context, id string, request RequestAboutDetailDTO, imageDetail *multipart.FileHeader) (*ResponseAboutDetailDTO, error) {
	errors, valid := request.ValidateAboutDetail()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	aboutDetail, err := s.aboutRepo.GetAboutDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("about detail tidak ditemukan: %v", err)
	}

	oldImageDetail := aboutDetail.ImageDetail

	var coverImageAboutDetailPath string
	if imageDetail != nil {
		coverImageAboutDetailPath, err = s.saveCoverImageAboutDetail(imageDetail)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan gambar baru: %v", err)
		}
	}

	aboutDetail.Description = request.Description
	if coverImageAboutDetailPath != "" {
		aboutDetail.ImageDetail = coverImageAboutDetailPath
	}

	updatedAboutDetail, err := s.aboutRepo.UpdateAboutDetail(ctx, id, aboutDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to update AboutDetail: %v", err)
	}

	if oldImageDetail != "" && coverImageAboutDetailPath != "" {
		if err := deleteCoverImageAbout(oldImageDetail); err != nil {
			log.Printf("Warning: failed to delete old detail image: %v", err)
		}
	}

	response, err := formatResponseAboutDetailDTO(updatedAboutDetail)
	if err != nil {
		return nil, fmt.Errorf("error formatting AboutDetail response: %v", err)
	}

	s.invalidateAboutDetailCaches(id, aboutDetail.AboutID)

	return response, nil
}

func (s *aboutService) DeleteAboutDetail(ctx context.Context, id string) error {
	aboutDetail, err := s.aboutRepo.GetAboutDetailByID(ctx, id)
	if err != nil {
		return fmt.Errorf("about detail tidak ditemukan: %v", err)
	}

	aboutID := aboutDetail.AboutID

	if aboutDetail.ImageDetail != "" {
		if err := deleteCoverImageAbout(aboutDetail.ImageDetail); err != nil {
			log.Printf("Warning: failed to delete detail image: %v", err)
		}
	}

	if err := s.aboutRepo.DeleteAboutDetail(ctx, id); err != nil {
		return fmt.Errorf("failed to delete AboutDetail: %v", err)
	}

	s.invalidateAboutDetailCaches(id, aboutID)

	return nil
}
