package worker
/* 
import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"rijig/config"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/model"
)

type CartWorker struct {
	cartService services.CartService
	cartRepo    repositories.CartRepository
	trashRepo   repositories.TrashRepository
}

func NewCartWorker(cartService services.CartService, cartRepo repositories.CartRepository, trashRepo repositories.TrashRepository) *CartWorker {
	return &CartWorker{
		cartService: cartService,
		cartRepo:    cartRepo,
		trashRepo:   trashRepo,
	}
}

func (w *CartWorker) AutoCommitExpiringCarts() error {
	ctx := context.Background()
	threshold := 1 * time.Minute

	keys, err := services.GetExpiringCartKeys(ctx, threshold)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	log.Printf("[CART-WORKER] Found %d carts expiring within 1 minute", len(keys))

	successCount := 0
	for _, key := range keys {
		userID := w.extractUserIDFromKey(key)
		if userID == "" {
			log.Printf("[CART-WORKER] Invalid key format: %s", key)
			continue
		}

		hasCart, err := w.cartRepo.HasExistingCart(ctx, userID)
		if err != nil {
			log.Printf("[CART-WORKER] Error checking existing cart for user %s: %v", userID, err)
			continue
		}

		if hasCart {

			if err := services.DeleteCartFromRedis(ctx, userID); err != nil {
				log.Printf("[CART-WORKER] Failed to delete Redis cache for user %s: %v", userID, err)
			} else {
				log.Printf("[CART-WORKER] Deleted Redis cache for user %s (already has DB cart)", userID)
			}
			continue
		}

		cartData, err := w.getCartFromRedis(ctx, key)
		if err != nil {
			log.Printf("[CART-WORKER] Failed to get cart data for key %s: %v", key, err)
			continue
		}

		if err := w.commitCartToDB(ctx, userID, cartData); err != nil {
			log.Printf("[CART-WORKER] Failed to commit cart for user %s: %v", userID, err)
			continue
		}

		if err := services.DeleteCartFromRedis(ctx, userID); err != nil {
			log.Printf("[CART-WORKER] Warning: Failed to delete Redis key after commit for user %s: %v", userID, err)
		}

		successCount++
		log.Printf("[CART-WORKER] Successfully auto-committed cart for user %s", userID)
	}

	log.Printf("[CART-WORKER] Auto-commit completed: %d successful commits", successCount)
	return nil
}

func (w *CartWorker) extractUserIDFromKey(key string) string {
	parts := strings.Split(key, ":")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}

func (w *CartWorker) getCartFromRedis(ctx context.Context, key string) (*dto.RequestCartDTO, error) {
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cart dto.RequestCartDTO
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (w *CartWorker) commitCartToDB(ctx context.Context, userID string, cartData *dto.RequestCartDTO) error {
	if len(cartData.CartItems) == 0 {
		return nil
	}

	totalAmount := 0.0
	totalPrice := 0.0
	var cartItems []model.CartItem

	for _, item := range cartData.CartItems {
		if item.Amount <= 0 {
			continue
		}

		trash, err := w.trashRepo.GetTrashCategoryByID(ctx, item.TrashID)
		if err != nil {
			log.Printf("[CART-WORKER] Warning: Skipping invalid trash category %s", item.TrashID)
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

	return w.cartRepo.CreateCartWithItems(ctx, newCart)
}
 */