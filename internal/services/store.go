package services

import (
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetStoreByID(storeID string) (dto.StoreResponseDTO, error) {
	store, err := repositories.GetStoreByID(storeID)
	if err != nil {
		return dto.StoreResponseDTO{}, err
	}

	return dto.StoreResponseDTO{
		ID:          store.ID,
		UserID:      store.UserID,
		StoreName:   store.StoreName,
		StoreLogo:   store.StoreLogo,
		StoreBanner: store.StoreBanner,
		StoreDesc:   store.StoreDesc,
		Follower:    store.Follower,
		StoreRating: store.StoreRating,
		CreatedAt:   utils.FormatDateToIndonesianFormat(store.CreatedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(store.UpdatedAt),
	}, nil
}

func GetStoresByUserID(userID string, limit, page int) ([]dto.StoreResponseDTO, error) {
	offset := (page - 1) * limit
	stores, err := repositories.GetStoresByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var storeResponses []dto.StoreResponseDTO
	for _, store := range stores {
		storeResponses = append(storeResponses, dto.StoreResponseDTO{
			ID:          store.ID,
			UserID:      store.UserID,
			StoreName:   store.StoreName,
			StoreLogo:   store.StoreLogo,
			StoreBanner: store.StoreBanner,
			StoreDesc:   store.StoreDesc,
			Follower:    store.Follower,
			StoreRating: store.StoreRating,
			CreatedAt:   utils.FormatDateToIndonesianFormat(store.CreatedAt),
			UpdatedAt:   utils.FormatDateToIndonesianFormat(store.UpdatedAt),
		})
	}

	return storeResponses, nil
}
