package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rijig/config"
	"rijig/dto"
)

var cartTTL = 30 * time.Minute

func getCartKey(userID string) string {
	return fmt.Sprintf("cart:user:%s", userID)
}

func SetCartToRedis(ctx context.Context, userID string, cart dto.RequestCartDTO) error {
	key := getCartKey(userID)

	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	err = config.RedisClient.Set(ctx, key, data, cartTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to save cart to redis: %w", err)
	}

	return nil
}

func GetCartFromRedis(ctx context.Context, userID string) (*dto.RequestCartDTO, error) {
	key := getCartKey(userID)
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cart dto.RequestCartDTO
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart data: %w", err)
	}

	return &cart, nil
}

func DeleteCartFromRedis(ctx context.Context, userID string) error {
	key := getCartKey(userID)
	return config.RedisClient.Del(ctx, key).Err()
}

func GetCartTTL(ctx context.Context, userID string) (time.Duration, error) {
	key := getCartKey(userID)
	return config.RedisClient.TTL(ctx, key).Result()
}

func UpdateOrAddCartItemToRedis(ctx context.Context, userID string, item dto.RequestCartItemDTO) error {
	cart, err := GetCartFromRedis(ctx, userID)
	if err != nil {

		cart = &dto.RequestCartDTO{
			CartItems: []dto.RequestCartItemDTO{item},
		}
		return SetCartToRedis(ctx, userID, *cart)
	}

	updated := false
	for i, ci := range cart.CartItems {
		if ci.TrashID == item.TrashID {
			cart.CartItems[i].Amount = item.Amount
			updated = true
			break
		}
	}
	if !updated {
		cart.CartItems = append(cart.CartItems, item)
	}

	return SetCartToRedis(ctx, userID, *cart)
}

func RemoveCartItemFromRedis(ctx context.Context, userID, trashID string) error {
	cart, err := GetCartFromRedis(ctx, userID)
	if err != nil {
		return err
	}

	updatedItems := make([]dto.RequestCartItemDTO, 0)
	for _, ci := range cart.CartItems {
		if ci.TrashID != trashID {
			updatedItems = append(updatedItems, ci)
		}
	}

	if len(updatedItems) == 0 {
		return DeleteCartFromRedis(ctx, userID)
	}

	cart.CartItems = updatedItems
	return SetCartToRedis(ctx, userID, *cart)
}
