package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type BannerService interface {
	CreateBanner(request dto.RequestBannerDTO, bannerImage *multipart.FileHeader) (*dto.ResponseBannerDTO, error)
	GetAllBanners() ([]dto.ResponseBannerDTO, error)
	GetBannerByID(id string) (*dto.ResponseBannerDTO, error)
	UpdateBanner(id string, request dto.RequestBannerDTO, bannerImage *multipart.FileHeader) (*dto.ResponseBannerDTO, error)
	DeleteBanner(id string) error
}

type bannerService struct {
	BannerRepo repositories.BannerRepository
}

func NewBannerService(bannerRepo repositories.BannerRepository) BannerService {
	return &bannerService{BannerRepo: bannerRepo}
}

func (s *bannerService) saveBannerImage(bannerImage *multipart.FileHeader) (string, error) {
	bannerImageDir := "./public/uploads/banners"
	if _, err := os.Stat(bannerImageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(bannerImageDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for banner image: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	extension := filepath.Ext(bannerImage.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	bannerImageFileName := fmt.Sprintf("%s_banner%s", uuid.New().String(), extension)
	bannerImagePath := filepath.Join(bannerImageDir, bannerImageFileName)

	src, err := bannerImage.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(bannerImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create banner image file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save banner image: %v", err)
	}

	return bannerImagePath, nil
}

func (s *bannerService) CreateBanner(request dto.RequestBannerDTO, bannerImage *multipart.FileHeader) (*dto.ResponseBannerDTO, error) {

	errors, valid := request.ValidateBannerInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	bannerImagePath, err := s.saveBannerImage(bannerImage)
	if err != nil {
		return nil, fmt.Errorf("failed to save banner image: %v", err)
	}

	banner := model.Banner{
		BannerName:  request.BannerName,
		BannerImage: bannerImagePath,
	}

	if err := s.BannerRepo.CreateBanner(&banner); err != nil {
		return nil, fmt.Errorf("failed to create banner: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(banner.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(banner.UpdatedAt)

	bannerResponseDTO := &dto.ResponseBannerDTO{
		ID:          banner.ID,
		BannerName:  banner.BannerName,
		BannerImage: banner.BannerImage,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	articlesCacheKey := "banners:all"
	err = utils.DeleteData(articlesCacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for all banners: %v\n", err)
	}

	cacheKey := fmt.Sprintf("banner:%s", banner.ID)
	cacheData := map[string]interface{}{
		"data": bannerResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching banner: %v\n", err)
	}

	banners, err := s.BannerRepo.FindAllBanners()
	if err == nil {
		var bannersDTO []dto.ResponseBannerDTO
		for _, b := range banners {
			createdAt, _ := utils.FormatDateToIndonesianFormat(b.CreatedAt)
			updatedAt, _ := utils.FormatDateToIndonesianFormat(b.UpdatedAt)

			bannersDTO = append(bannersDTO, dto.ResponseBannerDTO{
				ID:          b.ID,
				BannerName:  b.BannerName,
				BannerImage: b.BannerImage,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			})
		}

		cacheData = map[string]interface{}{
			"data": bannersDTO,
		}
		if err := utils.SetJSONData(articlesCacheKey, cacheData, time.Hour*24); err != nil {
			fmt.Printf("Error caching updated banners to Redis: %v\n", err)
		}
	} else {
		fmt.Printf("Error fetching all banners: %v\n", err)
	}

	return bannerResponseDTO, nil
}

func (s *bannerService) GetAllBanners() ([]dto.ResponseBannerDTO, error) {
	var banners []dto.ResponseBannerDTO

	cacheKey := "banners:all"
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {

		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				if bannerData, ok := item.(map[string]interface{}); ok {
					banners = append(banners, dto.ResponseBannerDTO{
						ID:          bannerData["id"].(string),
						BannerName:  bannerData["bannername"].(string),
						BannerImage: bannerData["bannerimage"].(string),
						CreatedAt:   bannerData["createdAt"].(string),
						UpdatedAt:   bannerData["updatedAt"].(string),
					})
				}
			}
			return banners, nil
		}
	}

	records, err := s.BannerRepo.FindAllBanners()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banners: %v", err)
	}

	for _, record := range records {
		createdAt, _ := utils.FormatDateToIndonesianFormat(record.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(record.UpdatedAt)

		banners = append(banners, dto.ResponseBannerDTO{
			ID:          record.ID,
			BannerName:  record.BannerName,
			BannerImage: record.BannerImage,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	cacheData := map[string]interface{}{
		"data": banners,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching banners: %v\n", err)
	}

	return banners, nil
}

func (s *bannerService) GetBannerByID(id string) (*dto.ResponseBannerDTO, error) {

	cacheKey := fmt.Sprintf("banner:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		if data, ok := cachedData["data"].(map[string]interface{}); ok {
			return &dto.ResponseBannerDTO{
				ID:          data["id"].(string),
				BannerName:  data["bannername"].(string),
				BannerImage: data["bannerimage"].(string),
				CreatedAt:   data["createdAt"].(string),
				UpdatedAt:   data["updatedAt"].(string),
			}, nil
		}
	}

	banner, err := s.BannerRepo.FindBannerByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banner by ID: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(banner.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(banner.UpdatedAt)

	bannerResponseDTO := &dto.ResponseBannerDTO{
		ID:          banner.ID,
		BannerName:  banner.BannerName,
		BannerImage: banner.BannerImage,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	cacheData := map[string]interface{}{
		"data": bannerResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching banner: %v\n", err)
	}

	return bannerResponseDTO, nil
}

func (s *bannerService) UpdateBanner(id string, request dto.RequestBannerDTO, bannerImage *multipart.FileHeader) (*dto.ResponseBannerDTO, error) {
	// Cari banner yang ingin diupdate
	banner, err := s.BannerRepo.FindBannerByID(id)
	if err != nil {
		return nil, fmt.Errorf("banner not found: %v", err)
	}

	// Update data banner
	banner.BannerName = request.BannerName
	if bannerImage != nil {
		// Hapus file lama jika ada gambar baru yang diupload
		bannerImagePath, err := s.saveBannerImage(bannerImage)
		if err != nil {
			return nil, fmt.Errorf("failed to save banner image: %v", err)
		}
		banner.BannerImage = bannerImagePath
	}

	// Simpan perubahan ke database
	if err := s.BannerRepo.UpdateBanner(id, banner); err != nil {
		return nil, fmt.Errorf("failed to update banner: %v", err)
	}

	// Format tanggal
	createdAt, _ := utils.FormatDateToIndonesianFormat(banner.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(banner.UpdatedAt)

	// Membuat Response DTO
	bannerResponseDTO := &dto.ResponseBannerDTO{
		ID:          banner.ID,
		BannerName:  banner.BannerName,
		BannerImage: banner.BannerImage,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Menghapus cache untuk banner yang lama
	cacheKey := fmt.Sprintf("banner:%s", id)
	err = utils.DeleteData(cacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for banner: %v\n", err)
	}

	// Cache banner yang terbaru
	cacheData := map[string]interface{}{
		"data": bannerResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching updated banner: %v\n", err)
	}

	// Menghapus dan memperbarui cache untuk seluruh banner
	articlesCacheKey := "banners:all"
	err = utils.DeleteData(articlesCacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for all banners: %v\n", err)
	}

	// Cache seluruh daftar banner yang terbaru
	banners, err := s.BannerRepo.FindAllBanners()
	if err == nil {
		var bannersDTO []dto.ResponseBannerDTO
		for _, b := range banners {
			createdAt, _ := utils.FormatDateToIndonesianFormat(b.CreatedAt)
			updatedAt, _ := utils.FormatDateToIndonesianFormat(b.UpdatedAt)

			bannersDTO = append(bannersDTO, dto.ResponseBannerDTO{
				ID:          b.ID,
				BannerName:  b.BannerName,
				BannerImage: b.BannerImage,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			})
		}

		cacheData = map[string]interface{}{
			"data": bannersDTO,
		}
		if err := utils.SetJSONData(articlesCacheKey, cacheData, time.Hour*24); err != nil {
			fmt.Printf("Error caching updated banners to Redis: %v\n", err)
		}
	} else {
		fmt.Printf("Error fetching all banners: %v\n", err)
	}

	return bannerResponseDTO, nil
}

// DeleteBanner - Menghapus banner dan memperbarui cache
func (s *bannerService) DeleteBanner(id string) error {
	// Hapus banner dari database
	if err := s.BannerRepo.DeleteBanner(id); err != nil {
		return fmt.Errorf("failed to delete banner: %v", err)
	}

	// Menghapus cache untuk banner yang dihapus
	cacheKey := fmt.Sprintf("banner:%s", id)
	err := utils.DeleteData(cacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for banner: %v\n", err)
	}

	// Menghapus cache untuk seluruh banner
	articlesCacheKey := "banners:all"
	err = utils.DeleteData(articlesCacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for all banners: %v\n", err)
	}

	return nil
}
