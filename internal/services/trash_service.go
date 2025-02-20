package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type TrashService interface {
	CreateCategory(request dto.RequestTrashCategoryDTO) (*dto.ResponseTrashCategoryDTO, error)
	AddDetailToCategory(request dto.RequestTrashDetailDTO) (*dto.ResponseTrashDetailDTO, error)

	GetCategories() ([]dto.ResponseTrashCategoryDTO, error)
	GetCategoryByID(id string) (*dto.ResponseTrashCategoryDTO, error)
	GetTrashDetailByID(id string) (*dto.ResponseTrashDetailDTO, error)

	UpdateCategory(id string, request dto.RequestTrashCategoryDTO) (*dto.ResponseTrashCategoryDTO, error)
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

func (s *trashService) CreateCategory(request dto.RequestTrashCategoryDTO) (*dto.ResponseTrashCategoryDTO, error) {
	errors, valid := request.ValidateTrashCategoryInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	category := model.TrashCategory{
		Name: request.Name,
	}

	if err := s.TrashRepo.CreateCategory(&category); err != nil {
		return nil, fmt.Errorf("failed to create category: %v", err)
	}

	categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := s.CacheCategoryAndDetails(category.ID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
		fmt.Printf("Error caching category: %v\n", err)
	}

	categories, err := s.TrashRepo.GetCategories()
	if err == nil {
		var categoriesDTO []dto.ResponseTrashCategoryDTO
		for _, c := range categories {
			categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
				ID:        c.ID,
				Name:      c.Name,
				CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
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

	detailResponseDTO := &dto.ResponseTrashDetailDTO{
		ID:          detail.ID,
		CategoryID:  detail.CategoryID,
		Description: detail.Description,
		Price:       detail.Price,
		CreatedAt:   detail.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   detail.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	cacheKey := fmt.Sprintf("detail:%s", detail.ID)
	cacheData := map[string]interface{}{
		"data": detailResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*6); err != nil {
		fmt.Printf("Error caching detail: %v\n", err)
	}

	category, err := s.TrashRepo.GetCategoryByID(detail.CategoryID)
	if err == nil {
		categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := s.CacheCategoryAndDetails(detail.CategoryID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
			fmt.Printf("Error caching updated category: %v\n", err)
		}
	} else {
		fmt.Printf("Error fetching category for cache update: %v\n", err)
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
				ID:        categoryData["id"].(string),
				Name:      categoryData["name"].(string),
				CreatedAt: categoryData["createdAt"].(string),
				UpdatedAt: categoryData["updatedAt"].(string),
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
		categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
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
		return &dto.ResponseTrashCategoryDTO{
			ID:        categoryData["id"].(string),
			Name:      categoryData["name"].(string),
			CreatedAt: categoryData["createdAt"].(string),
			UpdatedAt: categoryData["updatedAt"].(string),
			Details:   mapDetails(cachedCategory["details"]),
		}, nil
	}

	category, err := s.TrashRepo.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	categoryDTO := &dto.ResponseTrashCategoryDTO{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if category.Details != nil {
		var detailsDTO []dto.ResponseTrashDetailDTO
		for _, detail := range category.Details {
			detailsDTO = append(detailsDTO, dto.ResponseTrashDetailDTO{
				ID:          detail.ID,
				CategoryID:  detail.CategoryID,
				Description: detail.Description,
				Price:       detail.Price,
				CreatedAt:   detail.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:   detail.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		categoryDTO.Details = detailsDTO
	}

	if err := s.CacheCategoryAndDetails(category.ID, categoryDTO, categoryDTO.Details, time.Hour*6); err != nil {
		fmt.Printf("Error caching category and details: %v\n", err)
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

	detailDTO := &dto.ResponseTrashDetailDTO{
		ID:          detail.ID,
		CategoryID:  detail.CategoryID,
		Description: detail.Description,
		Price:       detail.Price,
		CreatedAt:   detail.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   detail.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	cacheData := map[string]interface{}{
		"data": detailDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching detail: %v\n", err)
	}

	return detailDTO, nil
}

func (s *trashService) UpdateCategory(id string, request dto.RequestTrashCategoryDTO) (*dto.ResponseTrashCategoryDTO, error) {
	errors, valid := request.ValidateTrashCategoryInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	if err := s.TrashRepo.UpdateCategoryName(id, request.Name); err != nil {
		return nil, fmt.Errorf("failed to update category: %v", err)
	}

	category, err := s.TrashRepo.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := s.CacheCategoryAndDetails(category.ID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
		fmt.Printf("Error caching updated category: %v\n", err)
	}

	allCategories, err := s.TrashRepo.GetCategories()
	if err == nil {
		var categoriesDTO []dto.ResponseTrashCategoryDTO
		for _, c := range allCategories {
			categoriesDTO = append(categoriesDTO, dto.ResponseTrashCategoryDTO{
				ID:        c.ID,
				Name:      c.Name,
				CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
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

	detailResponseDTO := &dto.ResponseTrashDetailDTO{
		ID:          detail.ID,
		CategoryID:  detail.CategoryID,
		Description: detail.Description,
		Price:       detail.Price,
		CreatedAt:   detail.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   detail.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	cacheKey := fmt.Sprintf("detail:%s", detail.ID)
	cacheData := map[string]interface{}{
		"data": detailResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*6); err != nil {
		fmt.Printf("Error caching updated detail: %v\n", err)
	}

	category, err := s.TrashRepo.GetCategoryByID(detail.CategoryID)
	if err == nil {
		categoryResponseDTO := &dto.ResponseTrashCategoryDTO{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := s.CacheCategoryAndDetails(detail.CategoryID, categoryResponseDTO, category.Details, time.Hour*6); err != nil {
			fmt.Printf("Error caching updated category: %v\n", err)
		}
	} else {
		fmt.Printf("Error fetching category for cache update: %v\n", err)
	}

	return detailResponseDTO, nil
}

func (s *trashService) DeleteCategory(id string) error {

	if err := s.TrashRepo.DeleteCategory(id); err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}

	cacheKey := fmt.Sprintf("category:%s", id)
	if err := utils.DeleteData(cacheKey); err != nil {
		fmt.Printf("Error clearing cache for deleted category: %v\n", err)
	}

	allCategoriesCacheKey := "categories:all"
	if err := utils.DeleteData(allCategoriesCacheKey); err != nil {
		fmt.Printf("Error clearing categories list cache: %v\n", err)
	}

	category, err := s.TrashRepo.GetCategoryByID(id)
	if err != nil {
		return fmt.Errorf("category not found after deletion: %v", err)
	}

	if category.Details != nil {
		for _, detail := range category.Details {
			detailCacheKey := fmt.Sprintf("detail:%s", detail.ID)

			if err := utils.DeleteData(detailCacheKey); err != nil {
				fmt.Printf("Error clearing cache for deleted detail: %v\n", err)
			}
		}
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

	cacheKey := fmt.Sprintf("detail:%s", id)
	if err := utils.DeleteData(cacheKey); err != nil {
		fmt.Printf("Error clearing cache for deleted detail: %v\n", err)
	}

	categoryCacheKey := fmt.Sprintf("category:%s", detail.CategoryID)
	if err := utils.DeleteData(categoryCacheKey); err != nil {
		fmt.Printf("Error clearing cache for category after detail deletion: %v\n", err)
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
