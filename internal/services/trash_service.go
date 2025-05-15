package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"

	"github.com/google/uuid"
)

type TrashService interface {
	CreateCategory(request dto.RequestTrashCategoryDTO, iconTrash *multipart.FileHeader) (*dto.ResponseTrashCategoryDTO, error)
	AddDetailToCategory(request dto.RequestTrashDetailDTO) (*dto.ResponseTrashDetailDTO, error)

	GetCategories() ([]dto.ResponseTrashCategoryDTO, error)
	GetCategoryByID(id string) (*dto.ResponseTrashCategoryDTO, error)
	GetTrashDetailByID(id string) (*dto.ResponseTrashDetailDTO, error)

	UpdateCategory(id string, request dto.RequestTrashCategoryDTO, iconPath *multipart.FileHeader) (*dto.ResponseTrashCategoryDTO, error)
	UpdateDetail(id string, request dto.RequestTrashDetailDTO) (*dto.ResponseTrashDetailDTO, error)

	DeleteCategory(id string) error
	DeleteDetail(id string) error
}

type trashService struct {
	TrashRepo repositories.TrashRepository
}

func NewTrashService(trashRepo repositories.TrashRepository) TrashService {
	return &trashService{TrashRepo: trashRepo}
}

func (s *trashService) saveIconOfTrash(iconTrash *multipart.FileHeader) (string, error) {
	pathImage := "/uploads/icontrash/"
	iconTrashDir := "./public" + os.Getenv("BASE_URL") + pathImage
	if _, err := os.Stat(iconTrashDir); os.IsNotExist(err) {

		if err := os.MkdirAll(iconTrashDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for icon trash: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := filepath.Ext(iconTrash.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	iconTrashFIleName := fmt.Sprintf("%s_icontrash%s", uuid.New().String(), extension)
	iconTrashPath := filepath.Join(iconTrashDir, iconTrashFIleName)

	src, err := iconTrash.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(iconTrashPath)
	if err != nil {
		return "", fmt.Errorf("failed to create icon trash file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save icon trash: %v", err)
	}

	iconTrashUrl := fmt.Sprintf("%s%s", pathImage, iconTrashFIleName)

	return iconTrashUrl, nil
}

func deleteIconTrashFIle(imagePath string) error {
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

func (s *trashService) CreateCategory(request dto.RequestTrashCategoryDTO, iconTrash *multipart.FileHeader) (*dto.ResponseTrashCategoryDTO, error) {

	parsedPrice, err := strconv.ParseFloat(request.EstimatedPrice, 64)
	fmt.Println("Received estimatedprice:", request.EstimatedPrice)
	if err != nil {
		return nil, fmt.Errorf("gagal memvalidasi harga: %v", err)
	} else {
		fmt.Printf("hasil parsing%v", parsedPrice)
	}

	icontrashPath, err := s.saveIconOfTrash(iconTrash)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan ikon sampah: %v", err)
	}

	category := model.TrashCategory{
		Name: request.Name,

		EstimatedPrice: parsedPrice,
		Icon:           icontrashPath,
	}

	if err := s.TrashRepo.CreateCategory(&category); err != nil {
		return nil, fmt.Errorf("failed to create category: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(category.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(category.UpdatedAt)

	categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
		ID:             category.ID,
		Name:           category.Name,
		EstimatedPrice: float64(category.EstimatedPrice),
		Icon:           category.Icon,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	if err := s.CacheCategoryAndDetails(category.ID, categoryResponseDTO, nil, time.Hour*6); err != nil {
		return nil, fmt.Errorf("error caching category: %v", err)
	}

	categories, err := s.TrashRepo.GetCategories()
	if err == nil {
		var categoriesDTO []dto.ResponseTrashCategoryDTO
		for _, c := range categories {
			ccreatedAt, _ := utils.FormatDateToIndonesianFormat(c.CreatedAt)
			cupdatedAt, _ := utils.FormatDateToIndonesianFormat(c.UpdatedAt)
			categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
				ID:             c.ID,
				Name:           c.Name,
				EstimatedPrice: float64(c.EstimatedPrice),
				Icon:           c.Icon,
				CreatedAt:      ccreatedAt,
				UpdatedAt:      cupdatedAt,
			})
		}

		if err := s.CacheCategoryList(categoriesDTO, time.Hour*6); err != nil {
			fmt.Printf("Error caching all categories: %v\n", err)
		}
	}

	return categoryResponseDTO, nil
}

func (s *trashService) AddDetailToCategory(request dto.RequestTrashDetailDTO) (*dto.ResponseTrashDetailDTO, error) {
	errors, valid := request.ValidateTrashDetailInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	detail := model.TrashDetail{
		CategoryID:  request.CategoryID,
		Description: request.Description,
		Price:       request.Price,
	}

	if err := s.TrashRepo.AddDetailToCategory(&detail); err != nil {
		return nil, fmt.Errorf("failed to add detail to category: %v", err)
	}

	dcreatedAt, _ := utils.FormatDateToIndonesianFormat(detail.CreatedAt)
	dupdatedAt, _ := utils.FormatDateToIndonesianFormat(detail.UpdatedAt)

	detailResponseDTO := &dto.ResponseTrashDetailDTO{
		ID:          detail.ID,
		CategoryID:  detail.CategoryID,
		Description: detail.Description,
		Price:       detail.Price,
		CreatedAt:   dcreatedAt,
		UpdatedAt:   dupdatedAt,
	}

	cacheKey := fmt.Sprintf("detail:%s", detail.ID)
	cacheData := map[string]interface{}{
		"data": detailResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*6); err != nil {
		return nil, fmt.Errorf("error caching detail: %v", err)
	}

	category, err := s.TrashRepo.GetCategoryByID(detail.CategoryID)

	if err == nil {

		ccreatedAt, _ := utils.FormatDateToIndonesianFormat(category.CreatedAt)
		cupdatedAt, _ := utils.FormatDateToIndonesianFormat(category.UpdatedAt)

		categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
			ID:        category.ID,
			Name:      category.Name,
			Icon:      category.Icon,
			CreatedAt: ccreatedAt,
			UpdatedAt: cupdatedAt,
		}

		if err := s.CacheCategoryAndDetails(detail.CategoryID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
			return nil, fmt.Errorf("error caching updated category: %v", err)
		}
	} else {
		return nil, fmt.Errorf("error fetching category for cache update: %v", err)
	}

	return detailResponseDTO, nil
}

func (s *trashService) GetCategories() ([]dto.ResponseTrashCategoryDTO, error) {
	cacheKey := "categories:all"
	cachedCategories, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedCategories != nil {
		var categoriesDTO []dto.ResponseTrashCategoryDTO
		for _, category := range cachedCategories["data"].([]interface{}) {
			categoryData := category.(map[string]interface{})
			categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
				ID:             categoryData["id"].(string),
				Name:           categoryData["name"].(string),
				EstimatedPrice: categoryData["estimatedprice"].(float64),
				Icon:           categoryData["icon"].(string),
				CreatedAt:      categoryData["createdAt"].(string),
				UpdatedAt:      categoryData["updatedAt"].(string),
			})
		}
		return categoriesDTO, nil
	}

	categories, err := s.TrashRepo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %v", err)
	}

	var categoriesDTO []dto.ResponseTrashCategoryDTO
	for _, category := range categories {

		createdAt, _ := utils.FormatDateToIndonesianFormat(category.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(category.UpdatedAt)
		categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
			ID:             category.ID,
			Name:           category.Name,
			EstimatedPrice: category.EstimatedPrice,
			Icon:           category.Icon,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		})
	}

	cacheData := map[string]interface{}{
		"data": categoriesDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*6); err != nil {
		fmt.Printf("Error caching categories: %v\n", err)
	}

	return categoriesDTO, nil
}

func (s *trashService) GetCategoryByID(id string) (*dto.ResponseTrashCategoryDTO, error) {
	cacheKey := fmt.Sprintf("category:%s", id)
	cachedCategory, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedCategory != nil {
		categoryData := cachedCategory["data"].(map[string]interface{})
		details := mapDetails(cachedCategory["details"])
		return &dto.ResponseTrashCategoryDTO{
			ID:             categoryData["id"].(string),
			Name:           categoryData["name"].(string),
			EstimatedPrice: categoryData["estimatedprice"].(float64),
			Icon:           categoryData["icon"].(string),
			CreatedAt:      categoryData["createdAt"].(string),
			UpdatedAt:      categoryData["updatedAt"].(string),
			Details:        details,
		}, nil
	}

	category, err := s.TrashRepo.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(category.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(category.UpdatedAt)

	categoryDTO := &dto.ResponseTrashCategoryDTO{
		ID:             category.ID,
		Name:           category.Name,
		EstimatedPrice: category.EstimatedPrice,
		Icon:           category.Icon,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	if category.Details != nil {
		var detailsDTO []dto.ResponseTrashDetailDTO
		for _, detail := range category.Details {
			createdAt, _ := utils.FormatDateToIndonesianFormat(detail.CreatedAt)
			updatedAt, _ := utils.FormatDateToIndonesianFormat(detail.UpdatedAt)
			detailsDTO = append(detailsDTO, dto.ResponseTrashDetailDTO{
				ID:          detail.ID,
				CategoryID:  detail.CategoryID,
				Description: detail.Description,
				Price:       detail.Price,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			})
		}
		categoryDTO.Details = detailsDTO
	}

	if err := s.CacheCategoryAndDetails(category.ID, categoryDTO, categoryDTO.Details, time.Hour*6); err != nil {
		return nil, fmt.Errorf("error caching category and details: %v", err)
	}

	return categoryDTO, nil
}

func (s *trashService) GetTrashDetailByID(id string) (*dto.ResponseTrashDetailDTO, error) {
	cacheKey := fmt.Sprintf("detail:%s", id)
	cachedDetail, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedDetail != nil {
		detailData := cachedDetail["data"].(map[string]interface{})
		return &dto.ResponseTrashDetailDTO{
			ID:          detailData["id"].(string),
			CategoryID:  detailData["category_id"].(string),
			Description: detailData["description"].(string),
			Price:       detailData["price"].(float64),
			CreatedAt:   detailData["createdAt"].(string),
			UpdatedAt:   detailData["updatedAt"].(string),
		}, nil
	}

	detail, err := s.TrashRepo.GetTrashDetailByID(id)
	if err != nil {
		return nil, fmt.Errorf("trash detail not found: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(detail.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(detail.UpdatedAt)

	detailDTO := &dto.ResponseTrashDetailDTO{
		ID:          detail.ID,
		CategoryID:  detail.CategoryID,
		Description: detail.Description,
		Price:       detail.Price,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	cacheData := map[string]interface{}{
		"data": detailDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		return nil, fmt.Errorf("error caching detail: %v", err)
	}

	return detailDTO, nil
}

func (s *trashService) UpdateCategory(id string, request dto.RequestTrashCategoryDTO, iconPath *multipart.FileHeader) (*dto.ResponseTrashCategoryDTO, error) {
	errors, valid := request.ValidateTrashCategoryInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	category, err := s.TrashRepo.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	if category.Icon != "" {
		err := deleteIconTrashFIle(category.Icon)
		if err != nil {
			return nil, fmt.Errorf("failed to delete old image: %v", err)
		}
	}

	var iconTrashPath string
	if iconPath != nil {
		iconTrashPath, err = s.saveIconOfTrash(iconPath)
		if err != nil {
			return nil, fmt.Errorf("failed to save card photo: %v", err)
		}
	}

	if iconTrashPath != "" {
		category.Icon = iconTrashPath
	}

	category, err = s.TrashRepo.UpdateCategory(id, category)
	if err != nil {
		log.Printf("Error updating trash category: %v", err)
		return nil, fmt.Errorf("failed to update category: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(category.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(category.UpdatedAt)

	categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
		ID:             category.ID,
		Name:           category.Name,
		EstimatedPrice: category.EstimatedPrice,
		Icon:           category.Icon,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	if err := s.CacheCategoryAndDetails(category.ID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
		return nil, fmt.Errorf("error caching updated category: %v", err)
	}

	allCategories, err := s.TrashRepo.GetCategories()
	if err == nil {
		var categoriesDTO []dto.ResponseTrashCategoryDTO
		for _, c := range allCategories {
			ccreatedAt, _ := utils.FormatDateToIndonesianFormat(c.CreatedAt)
			cupdatedAt, _ := utils.FormatDateToIndonesianFormat(c.UpdatedAt)
			categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
				ID:             c.ID,
				Name:           c.Name,
				EstimatedPrice: c.EstimatedPrice,
				Icon:           c.Icon,
				CreatedAt:      ccreatedAt,
				UpdatedAt:      cupdatedAt,
			})
		}

		if err := s.CacheCategoryList(categoriesDTO, time.Hour*6); err != nil {
			fmt.Printf("Error caching all categories: %v\n", err)
		}
	}

	return categoryResponseDTO, nil
}

func (s *trashService) UpdateDetail(id string, request dto.RequestTrashDetailDTO) (*dto.ResponseTrashDetailDTO, error) {
	errors, valid := request.ValidateTrashDetailInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	if err := s.TrashRepo.UpdateTrashDetail(id, request.Description, request.Price); err != nil {
		return nil, fmt.Errorf("failed to update detail: %v", err)
	}

	detail, err := s.TrashRepo.GetTrashDetailByID(id)
	if err != nil {
		return nil, fmt.Errorf("trash detail not found: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(detail.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(detail.UpdatedAt)

	detailResponseDTO := &dto.ResponseTrashDetailDTO{
		ID:          detail.ID,
		CategoryID:  detail.CategoryID,
		Description: detail.Description,
		Price:       detail.Price,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	cacheKey := fmt.Sprintf("detail:%s", detail.ID)
	cacheData := map[string]interface{}{
		"data": detailResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*6); err != nil {
		return nil, fmt.Errorf("error caching updated detail: %v", err)
	}

	category, err := s.TrashRepo.GetCategoryByID(detail.CategoryID)
	if err == nil {

		ccreatedAt, _ := utils.FormatDateToIndonesianFormat(category.CreatedAt)
		cupdatedAt, _ := utils.FormatDateToIndonesianFormat(category.UpdatedAt)

		categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
			ID:        category.ID,
			Name:      category.Name,
			Icon:      category.Icon,
			CreatedAt: ccreatedAt,
			UpdatedAt: cupdatedAt,
		}

		if err := s.CacheCategoryAndDetails(detail.CategoryID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
			return nil, fmt.Errorf("error caching updated category: %v", err)
		}
	} else {
		fmt.Printf("Error fetching category for cache update: %v\n", err)
	}

	return detailResponseDTO, nil
}

func (s *trashService) DeleteCategory(id string) error {
	detailsCacheKeyPrefix := "detail:"
	details, err := s.TrashRepo.GetDetailsByCategoryID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch details for category %s: %v", id, err)
	}

	for _, detail := range details {
		detailCacheKey := detailsCacheKeyPrefix + detail.ID
		if err := s.deleteCache(detailCacheKey); err != nil {
			return fmt.Errorf("error clearing cache for deleted detail: %v", err)
		}
	}

	category, err := s.TrashRepo.GetCategoryByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch category for deletion: %v", err)
	}

	if err := deleteIconTrashFIle(category.Icon); err != nil {
		return fmt.Errorf("error deleting icon for category %s: %v", id, err)
	}

	if err := s.TrashRepo.DeleteCategory(id); err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}

	if err := s.deleteCache("category:" + id); err != nil {
		return fmt.Errorf("error clearing cache for deleted category: %v", err)
	}

	if err := s.deleteCache("categories:all"); err != nil {
		return fmt.Errorf("error clearing categories list cache: %v", err)
	}

	return nil
}

func (s *trashService) DeleteDetail(id string) error {

	detail, err := s.TrashRepo.GetTrashDetailByID(id)
	if err != nil {
		return fmt.Errorf("trash detail not found: %v", err)
	}

	if err := s.TrashRepo.DeleteTrashDetail(id); err != nil {
		return fmt.Errorf("failed to delete detail: %v", err)
	}

	detailCacheKey := fmt.Sprintf("detail:%s", id)
	if err := s.deleteCache(detailCacheKey); err != nil {
		return fmt.Errorf("error clearing cache for deleted detail: %v", err)
	}

	categoryCacheKey := fmt.Sprintf("category:%s", detail.CategoryID)
	if err := s.deleteCache(categoryCacheKey); err != nil {
		return fmt.Errorf("error clearing cache for category after detail deletion: %v", err)
	}

	return nil
}

func mapDetails(details interface{}) []dto.ResponseTrashDetailDTO {
	var detailsDTO []dto.ResponseTrashDetailDTO
	if details != nil {
		for _, detail := range details.([]interface{}) {
			detailData := detail.(map[string]interface{})
			detailsDTO = append(detailsDTO, dto.ResponseTrashDetailDTO{
				ID:          detailData["id"].(string),
				CategoryID:  detailData["category_id"].(string),
				Description: detailData["description"].(string),
				Price:       detailData["price"].(float64),
				CreatedAt:   detailData["createdAt"].(string),
				UpdatedAt:   detailData["updatedAt"].(string),
			})
		}
	}
	return detailsDTO
}

func (s *trashService) CacheCategoryAndDetails(categoryID string, categoryData interface{}, detailsData interface{}, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("category:%s", categoryID)
	cacheData := map[string]interface{}{
		"data":    categoryData,
		"details": detailsData,
	}

	err := utils.SetJSONData(cacheKey, cacheData, expiration)
	if err != nil {
		return fmt.Errorf("error caching category and details: %v", err)
	}

	return nil
}

func (s *trashService) CacheCategoryList(categoriesData interface{}, expiration time.Duration) error {
	cacheKey := "categories:all"
	cacheData := map[string]interface{}{
		"data": categoriesData,
	}

	err := utils.SetJSONData(cacheKey, cacheData, expiration)
	if err != nil {
		return fmt.Errorf("error caching categories list: %v", err)
	}

	return nil
}

func (s *trashService) deleteCache(cacheKey string) error {
	if err := utils.DeleteData(cacheKey); err != nil {
		fmt.Printf("Error clearing cache for key: %v\n", cacheKey)
		return fmt.Errorf("error clearing cache for key %s: %v", cacheKey, err)
	}
	fmt.Printf("Deleted cache for key: %s\n", cacheKey)
	return nil
}
