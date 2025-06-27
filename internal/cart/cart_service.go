package cart

import (
	"context"
	"errors"
	"log"

	// "rijig/dto"
	// "rijig/internal/repositories"
	"rijig/internal/trash"
	"rijig/model"
)

type CartService interface {
	AddOrUpdateItem(ctx context.Context, userID string, req RequestCartItemDTO) error
	GetCart(ctx context.Context, userID string) (*ResponseCartDTO, error)
	DeleteItem(ctx context.Context, userID string, trashID string) error
	ClearCart(ctx context.Context, userID string) error
	Checkout(ctx context.Context, userID string) error
}

type cartService struct {
	repo      CartRepository
	trashRepo trash.TrashRepositoryInterface
}

func NewCartService(repo CartRepository, trashRepo trash.TrashRepositoryInterface) CartService {
	return &cartService{repo, trashRepo}
}

func (s *cartService) AddOrUpdateItem(ctx context.Context, userID string, req RequestCartItemDTO) error {
	if req.Amount <= 0 {
		return errors.New("amount harus lebih dari 0")
	}

	_, err := s.trashRepo.GetTrashCategoryByID(ctx, req.TrashID)
	if err != nil {
		return err
	}

	existingCart, err := GetCartFromRedis(ctx, userID)
	if err != nil {
		return err
	}

	if existingCart == nil {
		existingCart = &RequestCartDTO{
			CartItems: []RequestCartItemDTO{},
		}
	}

	updated := false
	for i, item := range existingCart.CartItems {
		if item.TrashID == req.TrashID {
			existingCart.CartItems[i].Amount = req.Amount
			updated = true
			break
		}
	}

	if !updated {
		existingCart.CartItems = append(existingCart.CartItems, RequestCartItemDTO{
			TrashID: req.TrashID,
			Amount:  req.Amount,
		})
	}

	return SetCartToRedis(ctx, userID, *existingCart)
}

func (s *cartService) GetCart(ctx context.Context, userID string) (*ResponseCartDTO, error) {

	cached, err := GetCartFromRedis(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cached != nil {

		if err := RefreshCartTTL(ctx, userID); err != nil {
			log.Printf("Warning: Failed to refresh cart TTL for user %s: %v", userID, err)
		}

		return s.buildResponseFromCache(ctx, userID, cached)
	}

	cart, err := s.repo.GetCartByUser(ctx, userID)
	if err != nil {

		return &ResponseCartDTO{
			ID:                  "",
			UserID:              userID,
			TotalAmount:         0,
			EstimatedTotalPrice: 0,
			CartItems:           []ResponseCartItemDTO{},
		}, nil

	}

	response := s.buildResponseFromDB(cart)

	cacheData := RequestCartDTO{CartItems: []RequestCartItemDTO{}}
	for _, item := range cart.CartItems {
		cacheData.CartItems = append(cacheData.CartItems, RequestCartItemDTO{
			TrashID: item.TrashCategoryID,
			Amount:  item.Amount,
		})
	}

	if err := SetCartToRedis(ctx, userID, cacheData); err != nil {
		log.Printf("Warning: Failed to cache cart for user %s: %v", userID, err)
	}

	return response, nil
}

func (s *cartService) DeleteItem(ctx context.Context, userID string, trashID string) error {

	existingCart, err := GetCartFromRedis(ctx, userID)
	if err != nil {
		return err
	}
	if existingCart == nil {
		return errors.New("keranjang tidak ditemukan")
	}

	filtered := []RequestCartItemDTO{}
	for _, item := range existingCart.CartItems {
		if item.TrashID != trashID {
			filtered = append(filtered, item)
		}
	}
	existingCart.CartItems = filtered

	return SetCartToRedis(ctx, userID, *existingCart)
}

func (s *cartService) ClearCart(ctx context.Context, userID string) error {

	if err := DeleteCartFromRedis(ctx, userID); err != nil {
		return err
	}

	return s.repo.DeleteCart(ctx, userID)
}

func (s *cartService) Checkout(ctx context.Context, userID string) error {

	cachedCart, err := GetCartFromRedis(ctx, userID)
	if err != nil {
		return err
	}

	if cachedCart != nil {
		if err := s.commitCartFromRedis(ctx, userID, cachedCart); err != nil {
			return err
		}
	}

	_, err = s.repo.GetCartByUser(ctx, userID)
	if err != nil {
		return err
	}

	DeleteCartFromRedis(ctx, userID)
	return s.repo.DeleteCart(ctx, userID)
}

func (s *cartService) buildResponseFromCache(ctx context.Context, userID string, cached *RequestCartDTO) (*ResponseCartDTO, error) {
	totalQty := 0.0
	totalPrice := 0.0
	items := []ResponseCartItemDTO{}

	for _, item := range cached.CartItems {
		trash, err := s.trashRepo.GetTrashCategoryByID(ctx, item.TrashID)
		if err != nil {
			log.Printf("Warning: Trash category %s not found for cached cart item", item.TrashID)
			continue
		}

		subtotal := item.Amount * trash.EstimatedPrice
		totalQty += item.Amount
		totalPrice += subtotal

		items = append(items, ResponseCartItemDTO{
			ID:                     "",
			TrashID:                item.TrashID,
			TrashName:              trash.Name,
			TrashIcon:              trash.IconTrash,
			TrashPrice:             trash.EstimatedPrice,
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: subtotal,
		})
	}

	return &ResponseCartDTO{
		ID:                  "-",
		UserID:              userID,
		TotalAmount:         totalQty,
		EstimatedTotalPrice: totalPrice,
		CartItems:           items,
	}, nil
}

func (s *cartService) buildResponseFromDB(cart *model.Cart) *ResponseCartDTO {
	var items []ResponseCartItemDTO
	for _, item := range cart.CartItems {
		items = append(items, ResponseCartItemDTO{
			ID:                     item.ID,
			TrashID:                item.TrashCategoryID,
			TrashName:              item.TrashCategory.Name,
			TrashIcon:              item.TrashCategory.IconTrash,
			TrashPrice:             item.TrashCategory.EstimatedPrice,
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: item.SubTotalEstimatedPrice,
		})
	}

	return &ResponseCartDTO{
		ID:                  cart.ID,
		UserID:              cart.UserID,
		TotalAmount:         cart.TotalAmount,
		EstimatedTotalPrice: cart.EstimatedTotalPrice,
		CartItems:           items,
	}
}

func (s *cartService) commitCartFromRedis(ctx context.Context, userID string, cachedCart *RequestCartDTO) error {
	if len(cachedCart.CartItems) == 0 {
		return nil
	}

	totalAmount := 0.0
	totalPrice := 0.0
	var cartItems []model.CartItem

	for _, item := range cachedCart.CartItems {
		trash, err := s.trashRepo.GetTrashCategoryByID(ctx, item.TrashID)
		if err != nil {
			log.Printf("Warning: Skipping invalid trash category %s during commit", item.TrashID)
			continue
		}

		subtotal := item.Amount * trash.EstimatedPrice
		totalAmount += item.Amount
		totalPrice += subtotal

		cartItems = append(cartItems, model.CartItem{
			TrashCategoryID:        item.TrashID,
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: subtotal,
		})
	}

	if len(cartItems) == 0 {
		return nil
	}

	newCart := &model.Cart{
		UserID:              userID,
		TotalAmount:         totalAmount,
		EstimatedTotalPrice: totalPrice,
		CartItems:           cartItems,
	}

	return s.repo.CreateCartWithItems(ctx, newCart)
}
