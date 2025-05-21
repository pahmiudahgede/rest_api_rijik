package services

import (
	"rijig/dto"
	"context"
)

type CartService interface {
	GetCart(ctx context.Context, userID string) (*dto.RequestCartDTO, error)
	AddOrUpdateItem(ctx context.Context, userID string, item dto.RequestCartItemDTO) error
	DeleteItem(ctx context.Context, userID string, trashID string) error
	ClearCart(ctx context.Context, userID string) error
}

type cartService struct{}

func NewCartService() CartService {
	return &cartService{}
}

func (s *cartService) GetCart(ctx context.Context, userID string) (*dto.RequestCartDTO, error) {
	return GetCartFromRedis(ctx, userID)
}

func (s *cartService) AddOrUpdateItem(ctx context.Context, userID string, item dto.RequestCartItemDTO) error {
	return UpdateOrAddCartItemToRedis(ctx, userID, item)
}

func (s *cartService) DeleteItem(ctx context.Context, userID string, trashID string) error {
	return RemoveCartItemFromRedis(ctx, userID, trashID)
}

func (s *cartService) ClearCart(ctx context.Context, userID string) error {
	return DeleteCartFromRedis(ctx, userID)
}