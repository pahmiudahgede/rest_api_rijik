package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProductsByStoreID(storeID string, limit, page int) ([]dto.ProductResponseWithSoldDTO, error) {

	offset := (page - 1) * limit

	products, err := repositories.GetProductsByStoreID(storeID, limit, offset)
	if err != nil {
		return nil, err
	}

	return mapProductsToDTO(products), nil
}

func GetProductsByUserID(userID string, limit, page int) ([]dto.ProductResponseWithSoldDTO, error) {
	offset := (page - 1) * limit
	products, err := repositories.GetProductsByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return mapProductsToDTO(products), nil
}

func mapProductsToDTO(products []domain.Product) []dto.ProductResponseWithSoldDTO {
	var productResponses []dto.ProductResponseWithSoldDTO
	for _, product := range products {
		var images []dto.ProductImageDTO
		for _, img := range product.ProductImages {
			images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
		}

		productResponses = append(productResponses, dto.ProductResponseWithSoldDTO{
			ID:            product.ID,
			StoreID:       product.StoreID,
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
	return productResponses
}

func GetProductByIDAndStoreID(productID, storeID string) (dto.ProductResponseWithSoldDTO, error) {
	product, err := repositories.GetProductByIDAndStoreID(productID, storeID)
	if err != nil {
		return dto.ProductResponseWithSoldDTO{}, err
	}

	var images []dto.ProductImageDTO
	for _, img := range product.ProductImages {
		images = append(images, dto.ProductImageDTO{ImageURL: img.ImageURL})
	}

	return dto.ProductResponseWithSoldDTO{
		ID:            product.ID,
		StoreID:       product.StoreID,
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

	trashDetail, err := repositories.GetTrashDetailByID(input.TrashDetailID)
	if err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	marketPrice := int64(trashDetail.Price)

	if err := dto.ValidateSalePrice(marketPrice, input.SalePrice); err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	product := &domain.Product{
		UserID:          userID,
		StoreID:         input.StoreID,
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

	trashDetail, err = repositories.GetTrashDetailByID(product.TrashDetailID)
	if err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	return dto.CreateProductResponseDTO{
		ID:            product.ID,
		StoreID:       product.StoreID,
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

func UpdateProduct(productID string, input dto.UpdateProductRequestDTO) (dto.CreateProductResponseDTO, error) {

	product, err := repositories.GetProductByID(productID)
	if err != nil {
		return dto.CreateProductResponseDTO{}, errors.New("product not found")
	}

	product.ProductTitle = input.ProductTitle
	product.TrashDetailID = input.TrashDetailID
	product.SalePrice = input.SalePrice
	product.Quantity = input.Quantity
	product.ProductDescribe = input.ProductDescribe

	var images []domain.ProductImage
	for _, imageURL := range input.ProductImages {
		images = append(images, domain.ProductImage{ImageURL: imageURL})
	}

	if err := repositories.UpdateProduct(&product, images); err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	trashDetail, err := repositories.GetTrashDetailByID(product.TrashDetailID)
	if err != nil {
		return dto.CreateProductResponseDTO{}, err
	}

	return dto.CreateProductResponseDTO{
		ID:            product.ID,
		StoreID:       product.StoreID,
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

func DeleteProduct(productID string) error {

	_, err := repositories.GetProductByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if err := repositories.DeleteProduct(productID); err != nil {
		return err
	}

	return nil
}
