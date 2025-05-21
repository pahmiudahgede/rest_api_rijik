package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"rijig/config"
	"rijig/dto"
	"rijig/model"
)

func CommitExpiredCartsToDB() error {
	ctx := context.Background()

	keys, err := config.RedisClient.Keys(ctx, "cart:user:*").Result()
	if err != nil {
		return fmt.Errorf("error fetching cart keys: %w", err)
	}

	for _, key := range keys {
		ttl, err := config.RedisClient.TTL(ctx, key).Result()
		if err != nil || ttl > 30*time.Second {
			continue
		}

		val, err := config.RedisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var cart dto.RequestCartDTO
		if err := json.Unmarshal([]byte(val), &cart); err != nil {
			continue
		}

		userID := extractUserIDFromKey(key)

		cartID := SaveCartToDB(ctx, userID, &cart)

		_ = config.RedisClient.Del(ctx, key).Err()

		fmt.Printf(
			"[AUTO-COMMIT] UserID: %s | CartID: %s | TotalItem: %d | EstimatedTotalPrice: %.2f | Committed at: %s\n",
			userID, cartID, len(cart.CartItems), calculateTotalEstimated(&cart), time.Now().Format(time.RFC3339),
		)

	}

	return nil
}

func extractUserIDFromKey(key string) string {

	parts := strings.Split(key, ":")
	if len(parts) == 3 {
		return parts[2]
	}
	return ""
}

func SaveCartToDB(ctx context.Context, userID string, cart *dto.RequestCartDTO) string {
	totalAmount := float32(0)
	totalPrice := float32(0)

	var cartItems []model.CartItem
	for _, item := range cart.CartItems {

		var trash model.TrashCategory
		if err := config.DB.First(&trash, "id = ?", item.TrashID).Error; err != nil {
			continue
		}

		subtotal := trash.EstimatedPrice * float64(item.Amount)
		totalAmount += item.Amount
		totalPrice += float32(subtotal)

		cartItems = append(cartItems, model.CartItem{
			TrashCategoryID:        item.TrashID,
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: float32(subtotal),
		})
	}

	newCart := model.Cart{
		UserID:              userID,
		TotalAmount:         totalAmount,
		EstimatedTotalPrice: totalPrice,
		CartItems:           cartItems,
	}

	if err := config.DB.WithContext(ctx).Create(&newCart).Error; err != nil {
		fmt.Printf("Error committing cart: %v\n", err)
	}

	return newCart.ID
}

func calculateTotalEstimated(cart *dto.RequestCartDTO) float32 {
	var total float32
	for _, item := range cart.CartItems {
		var trash model.TrashCategory
		if err := config.DB.First(&trash, "id = ?", item.TrashID).Error; err != nil {
			continue
		}
		total += item.Amount * float32(trash.EstimatedPrice)
	}
	return total
}
