package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type ProductService interface {
	SaveProductImage(file *multipart.FileHeader, imageType string) (string, error)
	CreateProduct(userID string, productDTO *dto.RequestProductDTO) (*dto.ResponseProductDTO, error)

	GetAllProductsByStoreID(userID string, page, limit int) ([]dto.ResponseProductDTO, int64, error)
	GetProductByID(productID string) (*dto.ResponseProductDTO, error)
}

type productService struct {
	productRepo repositories.ProductRepository
	storeRepo   repositories.StoreRepository
}

func NewProductService(productRepo repositories.ProductRepository, storeRepo repositories.StoreRepository) ProductService {
	return &productService{productRepo, storeRepo}
}

func (s *productService) CreateProduct(userID string, productDTO *dto.RequestProductDTO) (*dto.ResponseProductDTO, error) {
	store, err := s.storeRepo.FindStoreByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving store by user ID: %w", err)
	}
	if store == nil {
		return nil, fmt.Errorf("store not found for user %s", userID)
	}

	var imagePaths []string
	var productImages []model.ProductImage
	for _, file := range productDTO.ProductImages {
		imagePath, err := s.SaveProductImage(file, "product")
		if err != nil {
			return nil, fmt.Errorf("failed to save product image: %w", err)
		}
		imagePaths = append(imagePaths, imagePath)

		productImages = append(productImages, model.ProductImage{
			ImageURL: imagePath,
		})
	}

	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("at least one image is required for the product")
	}

	product := model.Product{
		StoreID:     store.ID,
		ProductName: productDTO.ProductName,
		Quantity:    productDTO.Quantity,
	}

	product.ProductImages = productImages

	if err := s.productRepo.CreateProduct(&product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	createdAt, err := utils.FormatDateToIndonesianFormat(product.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to format createdAt: %w", err)
	}
	updatedAt, err := utils.FormatDateToIndonesianFormat(product.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to format updatedAt: %w", err)
	}

	var productImagesDTO []dto.ResponseProductImageDTO
	for _, img := range product.ProductImages {
		productImagesDTO = append(productImagesDTO, dto.ResponseProductImageDTO{
			ID:        img.ID,
			ProductID: img.ProductID,
			ImageURL:  img.ImageURL,
		})
	}

	productDTOResponse := &dto.ResponseProductDTO{
		ID:            product.ID,
		StoreID:       product.StoreID,
		ProductName:   product.ProductName,
		Quantity:      product.Quantity,
		ProductImages: productImagesDTO,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return productDTOResponse, nil
}

func (s *productService) GetAllProductsByStoreID(userID string, page, limit int) ([]dto.ResponseProductDTO, int64, error) {

	store, err := s.storeRepo.FindStoreByUserID(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving store by user ID: %w", err)
	}
	if store == nil {
		return nil, 0, fmt.Errorf("store not found for user %s", userID)
	}

	total, err := s.productRepo.CountProductsByStoreID(store.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}

	products, err := s.productRepo.FindProductsByStoreID(store.ID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching products: %w", err)
	}

	var productDTOs []dto.ResponseProductDTO
	for _, product := range products {
		productImages, err := s.productRepo.FindProductImagesByProductID(product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("error fetching product images: %w", err)
		}

		var productImagesDTO []dto.ResponseProductImageDTO
		for _, img := range productImages {
			productImagesDTO = append(productImagesDTO, dto.ResponseProductImageDTO{
				ID:        img.ID,
				ProductID: img.ProductID,
				ImageURL:  img.ImageURL,
			})
		}

		createdAt, _ := utils.FormatDateToIndonesianFormat(product.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(product.UpdatedAt)

		productDTOs = append(productDTOs, dto.ResponseProductDTO{
			ID:            product.ID,
			StoreID:       product.StoreID,
			ProductName:   product.ProductName,
			Quantity:      product.Quantity,
			Saled:         product.Saled,
			ProductImages: productImagesDTO,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}

	return productDTOs, total, nil
}

func (s *productService) GetProductByID(productID string) (*dto.ResponseProductDTO, error) {

	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(product.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(product.UpdatedAt)

	productDTO := &dto.ResponseProductDTO{
		ID:          product.ID,
		StoreID:     product.StoreID,
		ProductName: product.ProductName,
		Quantity:    product.Quantity,
		Saled:       product.Saled,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	var productImagesDTO []dto.ResponseProductImageDTO
	for _, image := range product.ProductImages {
		productImagesDTO = append(productImagesDTO, dto.ResponseProductImageDTO{
			ID:        image.ID,
			ProductID: image.ProductID,
			ImageURL:  image.ImageURL,
		})
	}

	productDTO.ProductImages = productImagesDTO

	return productDTO, nil
}

func (s *productService) SaveProductImage(file *multipart.FileHeader, imageType string) (string, error) {
	imageDir := fmt.Sprintf("./public%s/uploads/store/%s", os.Getenv("BASE_URL"), imageType)
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for %s image: %v", imageType, err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	extension := filepath.Ext(file.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed for %s", imageType)
	}

	fileName := fmt.Sprintf("%s_%s%s", imageType, uuid.New().String(), extension)
	filePath := filepath.Join(imageDir, fileName)

	fileData, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer fileData.Close()

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create %s image file: %v", imageType, err)
	}
	defer outFile.Close()

	if _, err := outFile.ReadFrom(fileData); err != nil {
		return "", fmt.Errorf("failed to save %s image: %v", imageType, err)
	}

	return filepath.Join("/uploads/store/", imageType, fileName), nil
}
