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

	UpdateProduct(userID, productID string, productDTO *dto.RequestProductDTO) (*dto.ResponseProductDTO, error)
	DeleteProduct(productID string) error
	DeleteProducts(productIDs []string) error
	DeleteProductImage(imageID string) error
	DeleteProductImages(imageIDs []string) error
	deleteImageFile(imageID string) error
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

func (s *productService) UpdateProduct(userID, productID string, productDTO *dto.RequestProductDTO) (*dto.ResponseProductDTO, error) {

	store, err := s.storeRepo.FindStoreByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving store by user ID: %w", err)
	}
	if store == nil {
		return nil, fmt.Errorf("store not found for user %s", userID)
	}

	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product: %v", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	if product.StoreID != store.ID {
		return nil, fmt.Errorf("user does not own the store for this product")
	}

	if err := s.deleteProductImages(productID); err != nil {
		return nil, fmt.Errorf("failed to delete old product images: %v", err)
	}

	var productImages []model.ProductImage
	for _, file := range productDTO.ProductImages {
		imagePath, err := s.SaveProductImage(file, "product")
		if err != nil {
			return nil, fmt.Errorf("failed to save product image: %w", err)
		}

		productImages = append(productImages, model.ProductImage{
			ImageURL: imagePath,
		})
	}

	product.ProductName = productDTO.ProductName
	product.Quantity = productDTO.Quantity
	product.ProductImages = productImages

	if err := s.productRepo.UpdateProduct(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(product.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(product.UpdatedAt)

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

func (s *productService) DeleteProduct(productID string) error {

	if err := s.deleteProductImages(productID); err != nil {
		return fmt.Errorf("failed to delete associated product images: %w", err)
	}

	if err := s.productRepo.DeleteProduct(productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (s *productService) DeleteProducts(productIDs []string) error {

	for _, productID := range productIDs {
		if err := s.deleteProductImages(productID); err != nil {
			return fmt.Errorf("failed to delete associated images for product %s: %w", productID, err)
		}
	}

	if err := s.productRepo.DeleteProductsByID(productIDs); err != nil {
		return fmt.Errorf("failed to delete products: %w", err)
	}

	return nil
}

func (s *productService) DeleteProductImage(imageID string) error {

	if err := s.deleteImageFile(imageID); err != nil {
		return fmt.Errorf("failed to delete image file with ID %s: %w", imageID, err)
	}

	if err := s.productRepo.DeleteProductImageByID(imageID); err != nil {
		return fmt.Errorf("failed to delete product image from database: %w", err)
	}

	return nil
}

func (s *productService) DeleteProductImages(imageIDs []string) error {

	for _, imageID := range imageIDs {
		if err := s.deleteImageFile(imageID); err != nil {
			return fmt.Errorf("failed to delete image file with ID %s: %w", imageID, err)
		}
	}

	if err := s.productRepo.DeleteProductImagesByID(imageIDs); err != nil {
		return fmt.Errorf("failed to delete product images from database: %w", err)
	}

	return nil
}

func (s *productService) deleteProductImages(productID string) error {
	productImages, err := s.productRepo.FindProductImagesByProductID(productID)
	if err != nil {
		return fmt.Errorf("failed to fetch product images: %w", err)
	}

	for _, img := range productImages {
		if err := s.deleteImageFile(img.ID); err != nil {
			return fmt.Errorf("failed to delete image file: %w", err)
		}
	}

	if err := s.productRepo.DeleteProductImagesByProductID(productID); err != nil {
		return fmt.Errorf("failed to delete product images from database: %w", err)
	}

	return nil
}

func (s *productService) deleteImageFile(imageID string) error {
	productImage, err := s.productRepo.GetProductImageByID(imageID)
	if err != nil {
		return fmt.Errorf("failed to fetch product image: %w", err)
	}

	if productImage == nil {
		return fmt.Errorf("product image with ID %s not found", imageID)
	}

	baseURL := os.Getenv("BASE_URL")

	imagePath := fmt.Sprintf("./public%s%s", baseURL, productImage.ImageURL)

	if err := os.Remove(imagePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}
