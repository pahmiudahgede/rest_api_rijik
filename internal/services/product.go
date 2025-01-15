package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProductsByUserID(userID string, limit, page int) ([]dto.ProductResponseWithSoldDTO, error) {
	offset := (page - 1) * limit
	products, err := repositories.GetProductsByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponseWithSoldDTO
	for _, product := range products {
		var images []dto.ProductImageDTO
		for _, img := range product.ProductImages {
			images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
		}

		productResponses = append(productResponses, dto.ProductResponseWithSoldDTO{
			ID:            product.ID,
			UserID:        product.UserID,
			ProductTitle:  product.ProductTitle,
			ProductImages: images,
			TrashDetail: dto.TrashDetailResponseDTO{
				ID:          product.TrashDetail.ID,
				Description: product.TrashDetail.Description,
				Price:       product.TrashDetail.Price,
			},
			SalePrice:       product.SalePrice,
			Quantity:        product.Quantity,
			ProductDescribe: product.ProductDescribe,
			Sold:            product.Sold,
			CreatedAt:       utils.FormatDateToIndonesianFormat(product.CreatedAt),
			UpdatedAt:       utils.FormatDateToIndonesianFormat(product.UpdatedAt),
		})
	}

	return productResponses, nil
}

func GetProductByIDAndUserID(productID, userID string) (dto.ProductResponseWithSoldDTO, error) {
	product, err := repositories.GetProductByIDAndUserID(productID, userID)
	if err != nil {
		return dto.ProductResponseWithSoldDTO{}, err
	}

	var images []dto.ProductImageDTO
	for _, img := range product.ProductImages {
		images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
	}

	return dto.ProductResponseWithSoldDTO{
		ID:            product.ID,
		UserID:        product.UserID,
		ProductTitle:  product.ProductTitle,
		ProductImages: images,
		TrashDetail: dto.TrashDetailResponseDTO{
			ID:          product.TrashDetail.ID,
			Description: product.TrashDetail.Description,
			Price:       product.TrashDetail.Price,
		},
		SalePrice:       product.SalePrice,
		Quantity:        product.Quantity,
		ProductDescribe: product.ProductDescribe,
		Sold:            product.Sold,
		CreatedAt:       utils.FormatDateToIndonesianFormat(product.CreatedAt),
		UpdatedAt:       utils.FormatDateToIndonesianFormat(product.UpdatedAt),
	}, nil
}

func CreateProduct(input dto.CreateProductRequestDTO, userID string) (dto.CreateProductResponseDTO, error) {
	if err := dto.GetValidator().Struct(input); err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	product := &domain.Product{
		UserID:          userID,
		ProductTitle:    input.ProductTitle,
		TrashDetailID:   input.TrashDetailID,
		SalePrice:       input.SalePrice,
		Quantity:        input.Quantity,
		ProductDescribe: input.ProductDescribe,
	}

	var images []domain.ProductImage
	for _, imageURL := range input.ProductImages {
		images = append(images, domain.ProductImage{ImageURL: imageURL})
	}

	if err := repositories.CreateProduct(product, images); err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	trashDetail, err := repositories.GetTrashDetailByID(product.TrashDetailID)
	if err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	return dto.CreateProductResponseDTO{
		ID:            product.ID,
		UserID:        product.UserID,
		ProductTitle:  product.ProductTitle,
		ProductImages: input.ProductImages,
		TrashDetail: dto.TrashDetailResponseDTO{
			ID:          trashDetail.ID,
			Description: trashDetail.Description,
			Price:       trashDetail.Price,
		},
		SalePrice:       product.SalePrice,
		Quantity:        product.Quantity,
		ProductDescribe: product.ProductDescribe,
		CreatedAt:       utils.FormatDateToIndonesianFormat(product.CreatedAt),
		UpdatedAt:       utils.FormatDateToIndonesianFormat(product.UpdatedAt),
	}, nil
}

func UpdateProduct(productID, userID string, input dto.UpdateProductRequestDTO) (dto.ProductResponseDTO, error) {
	if err := dto.GetValidator().Struct(input); err != nil {
		return dto.ProductResponseDTO{}, err
	}

	product := &domain.Product{
		ID:              productID,
		UserID:          userID,
		ProductTitle:    input.ProductTitle,
		TrashDetailID:   input.TrashDetailID,
		SalePrice:       input.SalePrice,
		Quantity:        input.Quantity,
		ProductDescribe: input.ProductDescribe,
	}

	var images []domain.ProductImage
	for _, imageURL := range input.ProductImages {
		images = append(images, domain.ProductImage{ImageURL: imageURL})
	}

	if err := repositories.UpdateProduct(product, images); err != nil {
		return dto.ProductResponseDTO{}, err
	}

	updatedProduct, err := repositories.GetProductByID(productID)
	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	var productImages []dto.ProductImageDTO
	for _, img := range updatedProduct.ProductImages {
		productImages = append(productImages, dto.ProductImageDTO{ImageURL: img.ImageURL})
	}

	return dto.ProductResponseDTO{
		ID:            updatedProduct.ID,
		UserID:        updatedProduct.UserID,
		ProductTitle:  updatedProduct.ProductTitle,
		ProductImages: productImages,
		TrashDetail: dto.TrashDetailResponseDTO{
			ID:          updatedProduct.TrashDetail.ID,
			Description: updatedProduct.TrashDetail.Description,
			Price:       updatedProduct.TrashDetail.Price,
		},
		SalePrice:       updatedProduct.SalePrice,
		Quantity:        updatedProduct.Quantity,
		ProductDescribe: updatedProduct.ProductDescribe,
		CreatedAt:       utils.FormatDateToIndonesianFormat(updatedProduct.CreatedAt),
		UpdatedAt:       utils.FormatDateToIndonesianFormat(updatedProduct.UpdatedAt),
	}, nil
}

func DeleteProduct(productID, userID string) error {
	rowsAffected, err := repositories.DeleteProduct(productID, userID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found or not authorized to delete")
	}
	return nil
}
