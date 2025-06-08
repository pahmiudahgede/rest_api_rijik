package cart

import (
	"context"
	"fmt"
	"log"
	"time"

	"rijig/internal/trash"
	"rijig/model"
	"rijig/utils"

	"github.com/google/uuid"
)

type CartService struct {
	cartRepo  CartRepository
	trashRepo trash.TrashRepositoryInterface
}

func NewCartService(cartRepo CartRepository, trashRepo trash.TrashRepositoryInterface) *CartService {
	return &CartService{
		cartRepo:  cartRepo,
		trashRepo: trashRepo,
	}
}

func (s *CartService) AddToCart(ctx context.Context, userID, trashCategoryID string, amount float64) error {
	cartKey := fmt.Sprintf("cart:%s", userID)

	var cartItems map[string]model.CartItem
	err := utils.GetCache(cartKey, &cartItems)
	if err != nil && err.Error() != "ErrCacheMiss" {
		return fmt.Errorf("failed to get cart from cache: %w", err)
	}

	if cartItems == nil {
		cartItems = make(map[string]model.CartItem)
	}

	trashCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, trashCategoryID)
	if err != nil {
		return fmt.Errorf("failed to get trash category: %w", err)
	}

	cartItems[trashCategoryID] = model.CartItem{
		TrashCategoryID:        trashCategoryID,
		Amount:                 amount,
		SubTotalEstimatedPrice: amount * float64(trashCategory.EstimatedPrice),
	}

	return utils.SetCache(cartKey, cartItems, 24*time.Hour)
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID, trashCategoryID string) error {
	cartKey := fmt.Sprintf("cart:%s", userID)

	var cartItems map[string]model.CartItem
	err := utils.GetCache(cartKey, &cartItems)
	if err != nil {
		if err.Error() == "ErrCacheMiss" {
			return nil
		}
		return fmt.Errorf("failed to get cart from cache: %w", err)
	}

	delete(cartItems, trashCategoryID)

	if len(cartItems) == 0 {
		return utils.DeleteCache(cartKey)
	}

	return utils.SetCache(cartKey, cartItems, 24*time.Hour)
}

func (s *CartService) ClearCart(userID string) error {
	cartKey := fmt.Sprintf("cart:%s", userID)
	return utils.DeleteCache(cartKey)
}

func (s *CartService) GetCartFromRedis(ctx context.Context, userID string) (*CartResponse, error) {
	cartKey := fmt.Sprintf("cart:%s", userID)

	var cartItems map[string]model.CartItem
	err := utils.GetCache(cartKey, &cartItems)
	if err != nil {
		if err.Error() == "ErrCacheMiss" {
			return &CartResponse{
				ID:                  "N/A",
				UserID:              userID,
				TotalAmount:         0,
				EstimatedTotalPrice: 0,
				CartItems:           []CartItemResponse{},
			}, nil
		}
		return nil, fmt.Errorf("failed to get cart from cache: %w", err)
	}

	var totalAmount float64
	var estimatedTotal float64
	var cartItemDTOs []CartItemResponse

	for _, item := range cartItems {
		trashCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, item.TrashCategoryID)
		if err != nil {
			log.Printf("Failed to get trash category %s: %v", item.TrashCategoryID, err)
			continue
		}

		totalAmount += item.Amount
		estimatedTotal += item.SubTotalEstimatedPrice

		cartItemDTOs = append(cartItemDTOs, CartItemResponse{
			ID:                     uuid.NewString(),
			TrashID:                trashCategory.ID,
			TrashName:              trashCategory.Name,
			TrashIcon:              trashCategory.IconTrash,
			TrashPrice:             float64(trashCategory.EstimatedPrice),
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: item.SubTotalEstimatedPrice,
		})
	}

	resp := &CartResponse{
		ID:                  "N/A",
		UserID:              userID,
		TotalAmount:         totalAmount,
		EstimatedTotalPrice: estimatedTotal,
		CartItems:           cartItemDTOs,
	}

	return resp, nil
}

func (s *CartService) CommitCartToDatabase(ctx context.Context, userID string) error {
	cartKey := fmt.Sprintf("cart:%s", userID)

	var cartItems map[string]model.CartItem
	err := utils.GetCache(cartKey, &cartItems)
	if err != nil {
		if err.Error() == "ErrCacheMiss" {
			log.Printf("No cart items found in Redis for user: %s", userID)
			return fmt.Errorf("no cart items found")
		}
		return fmt.Errorf("failed to get cart from cache: %w", err)
	}

	if len(cartItems) == 0 {
		log.Printf("No items to commit for user: %s", userID)
		return fmt.Errorf("no items to commit")
	}

	hasCart, err := s.cartRepo.HasExistingCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing cart: %w", err)
	}

	var cart *model.Cart
	if hasCart {

		cart, err = s.cartRepo.GetCartByUser(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get existing cart: %w", err)
		}
	} else {

		cart, err = s.cartRepo.FindOrCreateCart(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
	}

	for _, item := range cartItems {
		trashCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, item.TrashCategoryID)
		if err != nil {
			log.Printf("Trash category not found for trashID: %s", item.TrashCategoryID)
			continue
		}

		err = s.cartRepo.AddOrUpdateCartItem(
			ctx,
			cart.ID,
			item.TrashCategoryID,
			item.Amount,
			float64(trashCategory.EstimatedPrice),
		)
		if err != nil {
			log.Printf("Failed to add/update cart item: %v", err)
			continue
		}
	}

	if err := s.cartRepo.UpdateCartTotals(ctx, cart.ID); err != nil {
		return fmt.Errorf("failed to update cart totals: %w", err)
	}

	if err := utils.DeleteCache(cartKey); err != nil {
		log.Printf("Failed to clear Redis cart: %v", err)
	}

	log.Printf("Cart committed successfully for user: %s", userID)
	return nil
}

func (s *CartService) GetCart(ctx context.Context, userID string) (*CartResponse, error) {

	cartRedis, err := s.GetCartFromRedis(ctx, userID)
	if err == nil && len(cartRedis.CartItems) > 0 {
		return cartRedis, nil
	}

	cartDB, err := s.cartRepo.GetCartByUser(ctx, userID)
	if err != nil {

		return &CartResponse{
			ID:                  "N/A",
			UserID:              userID,
			TotalAmount:         0,
			EstimatedTotalPrice: 0,
			CartItems:           []CartItemResponse{},
		}, nil
	}

	var items []CartItemResponse
	for _, item := range cartDB.CartItems {
		items = append(items, CartItemResponse{
			ID:                     item.ID,
			TrashID:                item.TrashCategoryID,
			TrashName:              item.TrashCategory.Name,
			TrashIcon:              item.TrashCategory.IconTrash,
			TrashPrice:             float64(item.TrashCategory.EstimatedPrice),
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: item.SubTotalEstimatedPrice,
		})
	}

	resp := &CartResponse{
		ID:                  cartDB.ID,
		UserID:              cartDB.UserID,
		TotalAmount:         cartDB.TotalAmount,
		EstimatedTotalPrice: cartDB.EstimatedTotalPrice,
		CartItems:           items,
	}

	return resp, nil
}

func (s *CartService) SyncCartFromDatabaseToRedis(ctx context.Context, userID string) error {

	cartDB, err := s.cartRepo.GetCartByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get cart from database: %w", err)
	}

	cartItems := make(map[string]model.CartItem)
	for _, item := range cartDB.CartItems {
		cartItems[item.TrashCategoryID] = model.CartItem{
			TrashCategoryID:        item.TrashCategoryID,
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: item.SubTotalEstimatedPrice,
		}
	}

	cartKey := fmt.Sprintf("cart:%s", userID)
	return utils.SetCache(cartKey, cartItems, 24*time.Hour)
}

func (s *CartService) GetCartItemCount(userID string) (int, error) {
	cartKey := fmt.Sprintf("cart:%s", userID)

	var cartItems map[string]model.CartItem
	err := utils.GetCache(cartKey, &cartItems)
	if err != nil {
		if err.Error() == "ErrCacheMiss" {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get cart from cache: %w", err)
	}

	return len(cartItems), nil
}

func (s *CartService) DeleteCart(ctx context.Context, userID string) error {

	cartKey := fmt.Sprintf("cart:%s", userID)
	if err := utils.DeleteCache(cartKey); err != nil {
		log.Printf("Failed to delete cart from Redis: %v", err)
	}

	return s.cartRepo.DeleteCart(ctx, userID)
}

func (s *CartService) UpdateCartWithDTO(ctx context.Context, userID string, cartDTO *RequestCartDTO) error {

	if errors, valid := cartDTO.ValidateRequestCartDTO(); !valid {
		return fmt.Errorf("validation failed: %v", errors)
	}

	cartKey := fmt.Sprintf("cart:%s", userID)
	cartItems := make(map[string]model.CartItem)

	for _, itemDTO := range cartDTO.CartItems {

		trashCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, itemDTO.TrashID)
		if err != nil {
			log.Printf("Failed to get trash category %s: %v", itemDTO.TrashID, err)
			continue
		}

		subtotal := itemDTO.Amount * float64(trashCategory.EstimatedPrice)

		cartItems[itemDTO.TrashID] = model.CartItem{
			TrashCategoryID:        itemDTO.TrashID,
			Amount:                 itemDTO.Amount,
			SubTotalEstimatedPrice: subtotal,
		}
	}

	return utils.SetCache(cartKey, cartItems, 24*time.Hour)
}

func (s *CartService) AddItemsToCart(ctx context.Context, userID string, items []RequestCartItemDTO) error {
	cartKey := fmt.Sprintf("cart:%s", userID)

	var cartItems map[string]model.CartItem
	err := utils.GetCache(cartKey, &cartItems)
	if err != nil && err.Error() != "ErrCacheMiss" {
		return fmt.Errorf("failed to get cart from cache: %w", err)
	}

	if cartItems == nil {
		cartItems = make(map[string]model.CartItem)
	}

	for _, itemDTO := range items {
		if itemDTO.TrashID == "" {
			continue
		}

		trashCategory, err := s.trashRepo.GetTrashCategoryByID(ctx, itemDTO.TrashID)
		if err != nil {
			log.Printf("Failed to get trash category %s: %v", itemDTO.TrashID, err)
			continue
		}

		subtotal := itemDTO.Amount * float64(trashCategory.EstimatedPrice)

		cartItems[itemDTO.TrashID] = model.CartItem{
			TrashCategoryID:        itemDTO.TrashID,
			Amount:                 itemDTO.Amount,
			SubTotalEstimatedPrice: subtotal,
		}
	}

	return utils.SetCache(cartKey, cartItems, 24*time.Hour)
}
