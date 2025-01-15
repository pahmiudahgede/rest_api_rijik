package services

import (
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProducts(limit, page int) ([]dto.ProductResponseDTO, error) {
	offset := (page - 1) * limit

	products, err := repositories.GetAllProducts(limit, offset)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponseDTO
	for _, product := range products {
		var images []dto.ProductImageDTO
		for _, img := range product.ProductImages {
			images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
		}

		trashDetail := dto.TrashDetailResponseDTO{
			ID:          product.TrashDetail.ID,
			Description: product.TrashDetail.Description,
			Price:       product.TrashDetail.Price,
		}

		productResponses = append(productResponses, dto.ProductResponseDTO{
			ID:              product.ID,
			UserID:          product.UserID,
			ProductTitle:    product.ProductTitle,
			ProductImages:   images,
			TrashDetail:     trashDetail,
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

func GetProductByID(productID string) (dto.ProductResponseDTO, error) {
	product, err := repositories.GetProductByID(productID)
	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	var images []dto.ProductImageDTO
	for _, img := range product.ProductImages {
		images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
	}

	trashDetail := dto.TrashDetailResponseDTO{
		ID:          product.TrashDetail.ID,
		Description: product.TrashDetail.Description,
		Price:       product.TrashDetail.Price,
	}

	productResponse := dto.ProductResponseDTO{
		ID:              product.ID,
		UserID:          product.UserID,
		ProductTitle:    product.ProductTitle,
		ProductImages:   images,
		TrashDetail:     trashDetail,
		SalePrice:       product.SalePrice,
		Quantity:        product.Quantity,
		ProductDescribe: product.ProductDescribe,
		Sold:            product.Sold,
		CreatedAt:       utils.FormatDateToIndonesianFormat(product.CreatedAt),
		UpdatedAt:       utils.FormatDateToIndonesianFormat(product.UpdatedAt),
	}

	return productResponse, nil
}
