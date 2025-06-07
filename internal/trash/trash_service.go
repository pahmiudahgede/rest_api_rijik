package trash

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rijig/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TrashServiceInterface interface {
	CreateTrashCategory(ctx context.Context, req RequestTrashCategoryDTO) (*ResponseTrashCategoryDTO, error)
	CreateTrashCategoryWithDetails(ctx context.Context, categoryReq RequestTrashCategoryDTO, detailsReq []RequestTrashDetailDTO) (*ResponseTrashCategoryDTO, error)
	CreateTrashCategoryWithIcon(ctx context.Context, req RequestTrashCategoryDTO, iconFile *multipart.FileHeader) (*ResponseTrashCategoryDTO, error)
	UpdateTrashCategory(ctx context.Context, id string, req RequestTrashCategoryDTO) (*ResponseTrashCategoryDTO, error)
	UpdateTrashCategoryWithIcon(ctx context.Context, id string, req RequestTrashCategoryDTO, iconFile *multipart.FileHeader) (*ResponseTrashCategoryDTO, error)
	GetAllTrashCategories(ctx context.Context) ([]ResponseTrashCategoryDTO, error)
	GetAllTrashCategoriesWithDetails(ctx context.Context) ([]ResponseTrashCategoryDTO, error)
	GetTrashCategoryByID(ctx context.Context, id string) (*ResponseTrashCategoryDTO, error)
	GetTrashCategoryByIDWithDetails(ctx context.Context, id string) (*ResponseTrashCategoryDTO, error)
	DeleteTrashCategory(ctx context.Context, id string) error

	CreateTrashDetail(ctx context.Context, req RequestTrashDetailDTO) (*ResponseTrashDetailDTO, error)
	CreateTrashDetailWithIcon(ctx context.Context, req RequestTrashDetailDTO, iconFile *multipart.FileHeader) (*ResponseTrashDetailDTO, error)
	AddTrashDetailToCategory(ctx context.Context, categoryID string, req RequestTrashDetailDTO) (*ResponseTrashDetailDTO, error)
	AddTrashDetailToCategoryWithIcon(ctx context.Context, categoryID string, req RequestTrashDetailDTO, iconFile *multipart.FileHeader) (*ResponseTrashDetailDTO, error)
	UpdateTrashDetail(ctx context.Context, id string, req RequestTrashDetailDTO) (*ResponseTrashDetailDTO, error)
	UpdateTrashDetailWithIcon(ctx context.Context, id string, req RequestTrashDetailDTO, iconFile *multipart.FileHeader) (*ResponseTrashDetailDTO, error)
	GetTrashDetailsByCategory(ctx context.Context, categoryID string) ([]ResponseTrashDetailDTO, error)
	GetTrashDetailByID(ctx context.Context, id string) (*ResponseTrashDetailDTO, error)
	DeleteTrashDetail(ctx context.Context, id string) error

	BulkCreateTrashDetails(ctx context.Context, categoryID string, detailsReq []RequestTrashDetailDTO) ([]ResponseTrashDetailDTO, error)
	BulkDeleteTrashDetails(ctx context.Context, detailIDs []string) error
	ReorderTrashDetails(ctx context.Context, categoryID string, orderedDetailIDs []string) error
}

type TrashService struct {
	trashRepo TrashRepositoryInterface
}

func NewTrashService(trashRepo TrashRepositoryInterface) TrashServiceInterface {
	return &TrashService{
		trashRepo: trashRepo,
	}
}

func (s *TrashService) saveIconOfTrash(iconTrash *multipart.FileHeader) (string, error) {
	pathImage := "/uploads/icontrash/"
	iconTrashDir := "./public" + os.Getenv("BASE_URL") + pathImage

	if _, err := os.Stat(iconTrashDir); os.IsNotExist(err) {
		if err := os.MkdirAll(iconTrashDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for icon trash: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := strings.ToLower(filepath.Ext(iconTrash.Filename))
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, .png, and .svg are allowed")
	}

	iconTrashFileName := fmt.Sprintf("%s_icontrash%s", uuid.New().String(), extension)
	iconTrashPath := filepath.Join(iconTrashDir, iconTrashFileName)

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

	iconTrashUrl := fmt.Sprintf("%s%s", pathImage, iconTrashFileName)
	return iconTrashUrl, nil
}

func (s *TrashService) saveIconOfTrashDetail(iconTrashDetail *multipart.FileHeader) (string, error) {
	pathImage := "/uploads/icontrashdetail/"
	iconTrashDetailDir := "./public" + os.Getenv("BASE_URL") + pathImage

	if _, err := os.Stat(iconTrashDetailDir); os.IsNotExist(err) {
		if err := os.MkdirAll(iconTrashDetailDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for icon trash detail: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := strings.ToLower(filepath.Ext(iconTrashDetail.Filename))
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, .png, and .svg are allowed")
	}

	iconTrashDetailFileName := fmt.Sprintf("%s_icontrashdetail%s", uuid.New().String(), extension)
	iconTrashDetailPath := filepath.Join(iconTrashDetailDir, iconTrashDetailFileName)

	src, err := iconTrashDetail.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(iconTrashDetailPath)
	if err != nil {
		return "", fmt.Errorf("failed to create icon trash detail file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save icon trash detail: %v", err)
	}

	iconTrashDetailUrl := fmt.Sprintf("%s%s", pathImage, iconTrashDetailFileName)
	return iconTrashDetailUrl, nil
}

func (s *TrashService) deleteIconTrashFile(imagePath string) error {
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

func (s *TrashService) deleteIconTrashDetailFile(imagePath string) error {
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

	log.Printf("Trash detail image deleted successfully: %s", absolutePath)
	return nil
}

func (s *TrashService) CreateTrashCategoryWithIcon(ctx context.Context, req RequestTrashCategoryDTO, iconFile *multipart.FileHeader) (*ResponseTrashCategoryDTO, error) {
	if errors, valid := req.ValidateRequestTrashCategoryDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	var iconUrl string
	var err error

	if iconFile != nil {
		iconUrl, err = s.saveIconOfTrash(iconFile)
		if err != nil {
			return nil, fmt.Errorf("failed to save icon: %w", err)
		}
	}

	category := &model.TrashCategory{
		Name:           req.Name,
		IconTrash:      iconUrl,
		EstimatedPrice: req.EstimatedPrice,
		Variety:        req.Variety,
	}

	if err := s.trashRepo.CreateTrashCategory(ctx, category); err != nil {

		if iconUrl != "" {
			s.deleteIconTrashFile(iconUrl)
		}
		return nil, fmt.Errorf("failed to create trash category: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTO(category)
	return response, nil
}

func (s *TrashService) UpdateTrashCategoryWithIcon(ctx context.Context, id string, req RequestTrashCategoryDTO, iconFile *multipart.FileHeader) (*ResponseTrashCategoryDTO, error) {
	if errors, valid := req.ValidateRequestTrashCategoryDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	existingCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing category: %w", err)
	}

	var iconUrl string = existingCategory.IconTrash

	if iconFile != nil {
		newIconUrl, err := s.saveIconOfTrash(iconFile)
		if err != nil {
			return nil, fmt.Errorf("failed to save new icon: %w", err)
		}
		iconUrl = newIconUrl
	}

	updates := map[string]interface{}{
		"name":            req.Name,
		"icon_trash":      iconUrl,
		"estimated_price": req.EstimatedPrice,
		"variety":         req.Variety,
	}

	if err := s.trashRepo.UpdateTrashCategory(ctx, id, updates); err != nil {

		if iconFile != nil && iconUrl != existingCategory.IconTrash {
			s.deleteIconTrashFile(iconUrl)
		}
		return nil, fmt.Errorf("failed to update trash category: %w", err)
	}

	if iconFile != nil && existingCategory.IconTrash != "" && iconUrl != existingCategory.IconTrash {
		if err := s.deleteIconTrashFile(existingCategory.IconTrash); err != nil {
			log.Printf("Warning: failed to delete old icon: %v", err)
		}
	}

	updatedCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated category: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTO(updatedCategory)
	return response, nil
}

func (s *TrashService) CreateTrashDetailWithIcon(ctx context.Context, req RequestTrashDetailDTO, iconFile *multipart.FileHeader) (*ResponseTrashDetailDTO, error) {
	if errors, valid := req.ValidateRequestTrashDetailDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	var iconUrl string
	var err error

	if iconFile != nil {
		iconUrl, err = s.saveIconOfTrashDetail(iconFile)
		if err != nil {
			return nil, fmt.Errorf("failed to save icon: %w", err)
		}
	}

	detail := &model.TrashDetail{
		TrashCategoryID: req.CategoryID,
		IconTrashDetail: iconUrl,
		Description:     req.Description,
		StepOrder:       req.StepOrder,
	}

	if err := s.trashRepo.CreateTrashDetail(ctx, detail); err != nil {

		if iconUrl != "" {
			s.deleteIconTrashDetailFile(iconUrl)
		}
		return nil, fmt.Errorf("failed to create trash detail: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(detail)
	return response, nil
}

func (s *TrashService) AddTrashDetailToCategoryWithIcon(ctx context.Context, categoryID string, req RequestTrashDetailDTO, iconFile *multipart.FileHeader) (*ResponseTrashDetailDTO, error) {
	if errors, valid := req.ValidateRequestTrashDetailDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	var iconUrl string
	var err error

	if iconFile != nil {
		iconUrl, err = s.saveIconOfTrashDetail(iconFile)
		if err != nil {
			return nil, fmt.Errorf("failed to save icon: %w", err)
		}
	}

	detail := &model.TrashDetail{
		IconTrashDetail: iconUrl,
		Description:     req.Description,
		StepOrder:       req.StepOrder,
	}

	if err := s.trashRepo.AddTrashDetailToCategory(ctx, categoryID, detail); err != nil {

		if iconUrl != "" {
			s.deleteIconTrashDetailFile(iconUrl)
		}
		return nil, fmt.Errorf("failed to add trash detail to category: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(detail)
	return response, nil
}

func (s *TrashService) UpdateTrashDetailWithIcon(ctx context.Context, id string, req RequestTrashDetailDTO, iconFile *multipart.FileHeader) (*ResponseTrashDetailDTO, error) {
	if errors, valid := req.ValidateRequestTrashDetailDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	existingDetail, err := s.trashRepo.GetTrashDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing detail: %w", err)
	}

	var iconUrl string = existingDetail.IconTrashDetail

	if iconFile != nil {
		newIconUrl, err := s.saveIconOfTrashDetail(iconFile)
		if err != nil {
			return nil, fmt.Errorf("failed to save new icon: %w", err)
		}
		iconUrl = newIconUrl
	}

	updates := map[string]interface{}{
		"icon_trash_detail": iconUrl,
		"description":       req.Description,
		"step_order":        req.StepOrder,
	}

	if err := s.trashRepo.UpdateTrashDetail(ctx, id, updates); err != nil {

		if iconFile != nil && iconUrl != existingDetail.IconTrashDetail {
			s.deleteIconTrashDetailFile(iconUrl)
		}
		return nil, fmt.Errorf("failed to update trash detail: %w", err)
	}

	if iconFile != nil && existingDetail.IconTrashDetail != "" && iconUrl != existingDetail.IconTrashDetail {
		if err := s.deleteIconTrashDetailFile(existingDetail.IconTrashDetail); err != nil {
			log.Printf("Warning: failed to delete old icon: %v", err)
		}
	}

	updatedDetail, err := s.trashRepo.GetTrashDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated detail: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(updatedDetail)
	return response, nil
}

func (s *TrashService) CreateTrashCategory(ctx context.Context, req RequestTrashCategoryDTO) (*ResponseTrashCategoryDTO, error) {
	if errors, valid := req.ValidateRequestTrashCategoryDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	category := &model.TrashCategory{
		Name:           req.Name,
		IconTrash:      req.IconTrash,
		EstimatedPrice: req.EstimatedPrice,
		Variety:        req.Variety,
	}

	if err := s.trashRepo.CreateTrashCategory(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create trash category: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTO(category)
	return response, nil
}

func (s *TrashService) CreateTrashCategoryWithDetails(ctx context.Context, categoryReq RequestTrashCategoryDTO, detailsReq []RequestTrashDetailDTO) (*ResponseTrashCategoryDTO, error) {
	if errors, valid := categoryReq.ValidateRequestTrashCategoryDTO(); !valid {
		return nil, fmt.Errorf("category validation failed: %v", errors)
	}

	for i, detailReq := range detailsReq {
		if errors, valid := detailReq.ValidateRequestTrashDetailDTO(); !valid {
			return nil, fmt.Errorf("detail %d validation failed: %v", i+1, errors)
		}
	}

	category := &model.TrashCategory{
		Name:           categoryReq.Name,
		IconTrash:      categoryReq.IconTrash,
		EstimatedPrice: categoryReq.EstimatedPrice,
		Variety:        categoryReq.Variety,
	}

	details := make([]model.TrashDetail, len(detailsReq))
	for i, detailReq := range detailsReq {
		details[i] = model.TrashDetail{
			IconTrashDetail: detailReq.IconTrashDetail,
			Description:     detailReq.Description,
			StepOrder:       detailReq.StepOrder,
		}
	}

	if err := s.trashRepo.CreateTrashCategoryWithDetails(ctx, category, details); err != nil {
		return nil, fmt.Errorf("failed to create trash category with details: %w", err)
	}

	createdCategory, err := s.trashRepo.GetTrashCategoryByIDWithDetails(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created category: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTOWithDetails(createdCategory)
	return response, nil
}

func (s *TrashService) UpdateTrashCategory(ctx context.Context, id string, req RequestTrashCategoryDTO) (*ResponseTrashCategoryDTO, error) {
	if errors, valid := req.ValidateRequestTrashCategoryDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	exists, err := s.trashRepo.CheckTrashCategoryExists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}
	if !exists {
		return nil, errors.New("trash category not found")
	}

	updates := map[string]interface{}{
		"name":            req.Name,
		"icon_trash":      req.IconTrash,
		"estimated_price": req.EstimatedPrice,
		"variety":         req.Variety,
	}

	if err := s.trashRepo.UpdateTrashCategory(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update trash category: %w", err)
	}

	updatedCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated category: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTO(updatedCategory)
	return response, nil
}

func (s *TrashService) GetAllTrashCategories(ctx context.Context) ([]ResponseTrashCategoryDTO, error) {
	categories, err := s.trashRepo.GetAllTrashCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash categories: %w", err)
	}

	responses := make([]ResponseTrashCategoryDTO, len(categories))
	for i, category := range categories {
		responses[i] = *s.convertTrashCategoryToResponseDTO(&category)
	}

	return responses, nil
}

func (s *TrashService) GetAllTrashCategoriesWithDetails(ctx context.Context) ([]ResponseTrashCategoryDTO, error) {
	categories, err := s.trashRepo.GetAllTrashCategoriesWithDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash categories with details: %w", err)
	}

	responses := make([]ResponseTrashCategoryDTO, len(categories))
	for i, category := range categories {
		responses[i] = *s.convertTrashCategoryToResponseDTOWithDetails(&category)
	}

	return responses, nil
}

func (s *TrashService) GetTrashCategoryByID(ctx context.Context, id string) (*ResponseTrashCategoryDTO, error) {
	category, err := s.trashRepo.GetTrashCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash category: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTO(category)
	return response, nil
}

func (s *TrashService) GetTrashCategoryByIDWithDetails(ctx context.Context, id string) (*ResponseTrashCategoryDTO, error) {
	category, err := s.trashRepo.GetTrashCategoryByIDWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash category with details: %w", err)
	}

	response := s.convertTrashCategoryToResponseDTOWithDetails(category)
	return response, nil
}

func (s *TrashService) DeleteTrashCategory(ctx context.Context, id string) error {

	category, err := s.trashRepo.GetTrashCategoryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	if err := s.trashRepo.DeleteTrashCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to delete trash category: %w", err)
	}

	if category.IconTrash != "" {
		if err := s.deleteIconTrashFile(category.IconTrash); err != nil {
			log.Printf("Warning: failed to delete category icon: %v", err)
		}
	}

	return nil
}

func (s *TrashService) CreateTrashDetail(ctx context.Context, req RequestTrashDetailDTO) (*ResponseTrashDetailDTO, error) {
	if errors, valid := req.ValidateRequestTrashDetailDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	detail := &model.TrashDetail{
		TrashCategoryID: req.CategoryID,
		IconTrashDetail: req.IconTrashDetail,
		Description:     req.Description,
		StepOrder:       req.StepOrder,
	}

	if err := s.trashRepo.CreateTrashDetail(ctx, detail); err != nil {
		return nil, fmt.Errorf("failed to create trash detail: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(detail)
	return response, nil
}

func (s *TrashService) AddTrashDetailToCategory(ctx context.Context, categoryID string, req RequestTrashDetailDTO) (*ResponseTrashDetailDTO, error) {
	if errors, valid := req.ValidateRequestTrashDetailDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	detail := &model.TrashDetail{
		IconTrashDetail: req.IconTrashDetail,
		Description:     req.Description,
		StepOrder:       req.StepOrder,
	}

	if err := s.trashRepo.AddTrashDetailToCategory(ctx, categoryID, detail); err != nil {
		return nil, fmt.Errorf("failed to add trash detail to category: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(detail)
	return response, nil
}

func (s *TrashService) UpdateTrashDetail(ctx context.Context, id string, req RequestTrashDetailDTO) (*ResponseTrashDetailDTO, error) {
	if errors, valid := req.ValidateRequestTrashDetailDTO(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	exists, err := s.trashRepo.CheckTrashDetailExists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check detail existence: %w", err)
	}
	if !exists {
		return nil, errors.New("trash detail not found")
	}

	updates := map[string]interface{}{
		"icon_trash_detail": req.IconTrashDetail,
		"description":       req.Description,
		"step_order":        req.StepOrder,
	}

	if err := s.trashRepo.UpdateTrashDetail(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update trash detail: %w", err)
	}

	updatedDetail, err := s.trashRepo.GetTrashDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated detail: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(updatedDetail)
	return response, nil
}

func (s *TrashService) GetTrashDetailsByCategory(ctx context.Context, categoryID string) ([]ResponseTrashDetailDTO, error) {
	exists, err := s.trashRepo.CheckTrashCategoryExists(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}
	if !exists {
		return nil, errors.New("trash category not found")
	}

	details, err := s.trashRepo.GetTrashDetailsByCategory(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash details: %w", err)
	}

	responses := make([]ResponseTrashDetailDTO, len(details))
	for i, detail := range details {
		responses[i] = *s.convertTrashDetailToResponseDTO(&detail)
	}

	return responses, nil
}

func (s *TrashService) GetTrashDetailByID(ctx context.Context, id string) (*ResponseTrashDetailDTO, error) {
	detail, err := s.trashRepo.GetTrashDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash detail: %w", err)
	}

	response := s.convertTrashDetailToResponseDTO(detail)
	return response, nil
}

func (s *TrashService) DeleteTrashDetail(ctx context.Context, id string) error {

	detail, err := s.trashRepo.GetTrashDetailByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get detail: %w", err)
	}

	if err := s.trashRepo.DeleteTrashDetail(ctx, id); err != nil {
		return fmt.Errorf("failed to delete trash detail: %w", err)
	}

	if detail.IconTrashDetail != "" {
		if err := s.deleteIconTrashDetailFile(detail.IconTrashDetail); err != nil {
			log.Printf("Warning: failed to delete detail icon: %v", err)
		}
	}

	return nil
}

func (s *TrashService) BulkCreateTrashDetails(ctx context.Context, categoryID string, detailsReq []RequestTrashDetailDTO) ([]ResponseTrashDetailDTO, error) {
	exists, err := s.trashRepo.CheckTrashCategoryExists(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}
	if !exists {
		return nil, errors.New("trash category not found")
	}

	for i, detailReq := range detailsReq {
		if errors, valid := detailReq.ValidateRequestTrashDetailDTO(); !valid {
			return nil, fmt.Errorf("detail %d validation failed: %v", i+1, errors)
		}
	}

	responses := make([]ResponseTrashDetailDTO, len(detailsReq))
	for i, detailReq := range detailsReq {
		response, err := s.AddTrashDetailToCategory(ctx, categoryID, detailReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create detail %d: %w", i+1, err)
		}
		responses[i] = *response
	}

	return responses, nil
}

func (s *TrashService) BulkDeleteTrashDetails(ctx context.Context, detailIDs []string) error {
	for _, id := range detailIDs {
		if err := s.DeleteTrashDetail(ctx, id); err != nil {
			return fmt.Errorf("failed to delete detail %s: %w", id, err)
		}
	}
	return nil
}

func (s *TrashService) ReorderTrashDetails(ctx context.Context, categoryID string, orderedDetailIDs []string) error {
	exists, err := s.trashRepo.CheckTrashCategoryExists(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to check category existence: %w", err)
	}
	if !exists {
		return errors.New("trash category not found")
	}

	for i, detailID := range orderedDetailIDs {
		updates := map[string]interface{}{
			"step_order": i + 1,
		}
		if err := s.trashRepo.UpdateTrashDetail(ctx, detailID, updates); err != nil {
			return fmt.Errorf("failed to reorder detail %s: %w", detailID, err)
		}
	}

	return nil
}

func (s *TrashService) convertTrashCategoryToResponseDTO(category *model.TrashCategory) *ResponseTrashCategoryDTO {
	return &ResponseTrashCategoryDTO{
		ID:             category.ID,
		TrashName:      category.Name,
		TrashIcon:      category.IconTrash,
		EstimatedPrice: category.EstimatedPrice,
		Variety:        category.Variety,
		CreatedAt:      category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      category.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *TrashService) convertTrashCategoryToResponseDTOWithDetails(category *model.TrashCategory) *ResponseTrashCategoryDTO {
	response := s.convertTrashCategoryToResponseDTO(category)

	details := make([]ResponseTrashDetailDTO, len(category.Details))
	for i, detail := range category.Details {
		details[i] = *s.convertTrashDetailToResponseDTO(&detail)
	}
	response.TrashDetail = details

	return response
}

func (s *TrashService) convertTrashDetailToResponseDTO(detail *model.TrashDetail) *ResponseTrashDetailDTO {
	return &ResponseTrashDetailDTO{
		ID:              detail.ID,
		CategoryID:      detail.TrashCategoryID,
		IconTrashDetail: detail.IconTrashDetail,
		Description:     detail.Description,
		StepOrder:       detail.StepOrder,
		CreatedAt:       detail.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       detail.UpdatedAt.Format(time.RFC3339),
	}
}
