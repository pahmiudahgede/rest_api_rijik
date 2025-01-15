package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProductsByUserID(userID string, limit, page int) ([]dto.ProductResponseDTO, error) {
	offset := (page - 1) * limit
	products, err := repositories.GetProductsByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponseDTO
	for _, product := range products {
		var images []dto.ProductImageDTO
		for _, img := range product.ProductImages {
			images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
		}

		productResponses = append(productResponses, dto.ProductResponseDTO{
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

func GetProductByIDAndUserID(productID, userID string) (dto.ProductResponseDTO, error) {
	product, err := repositories.GetProductByIDAndUserID(productID, userID)
	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	var images []dto.ProductImageDTO
	for _, img := range product.ProductImages {
		images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
	}

	return dto.ProductResponseDTO{
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

	return dto.CreateProductResponseDTO{
		ID:              product.ID,
		UserID:          product.UserID,
		ProductTitle:    product.ProductTitle,
		ProductImages:   input.ProductImages,
		TrashDetailID:   product.TrashDetailID,
		SalePrice:       product.SalePrice,
		Quantity:        product.Quantity,
		ProductDescribe: product.ProductDescribe,
	}, nil
}
