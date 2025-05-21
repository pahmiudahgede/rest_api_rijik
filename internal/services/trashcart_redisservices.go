package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"rijig/config"
	"rijig/dto"

	"github.com/go-redis/redis/v8"
)

const cartTTL = 1 * time.Minute

func getCartKey(userID string) string {
	return fmt.Sprintf("cart:%s", userID)
}

func GetCartItems(userID string) ([]dto.RequestCartItems, error) {
	key := getCartKey(userID)
	val, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var items []dto.RequestCartItems
	err = json.Unmarshal([]byte(val), &items)
	if err != nil {
		return nil, err
	}
	
	return items, nil
}

func AddOrUpdateCartItem(userID string, newItem dto.RequestCartItems) error {
	key := getCartKey(userID)
	var cartItems []dto.RequestCartItems

	val, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err == nil && val != "" {
		json.Unmarshal([]byte(val), &cartItems)
	}

	updated := false
	for i, item := range cartItems {
		if item.TrashID == newItem.TrashID {
			if newItem.Amount == 0 {
				cartItems = append(cartItems[:i], cartItems[i+1:]...)
			} else {
				cartItems[i].Amount = newItem.Amount
			}
			updated = true
			break
		}
	}

	if !updated && newItem.Amount > 0 {
		cartItems = append(cartItems, newItem)
	}

	return setCartItems(key, cartItems)
}

func DeleteCartItem(userID, trashID string) error {
	key := fmt.Sprintf("cart:%s", userID)
	items, err := GetCartItems(userID)

	if err == redis.Nil {

		log.Printf("No cart found in Redis for user: %s", userID)
		return fmt.Errorf("no cart found")
	}

	if err != nil {
		log.Printf("Redis error: %v", err)
		return err
	}

	index := -1
	for i, item := range items {
		if item.TrashID == trashID {
			index = i
			break
		}
	}

	if index == -1 {
		log.Printf("TrashID %s not found in cart for user %s", trashID, userID)
		return fmt.Errorf("trashid not found")
	}

	items = append(items[:index], items[index+1:]...)

	if len(items) == 0 {
		return config.RedisClient.Del(config.Ctx, key).Err()
	}

	return setCartItems(key, items)
}

func ClearCart(userID string) error {
	key := getCartKey(userID)
	return config.RedisClient.Del(config.Ctx, key).Err()
}

func RefreshCartTTL(userID string) error {
	key := getCartKey(userID)
	return config.RedisClient.Expire(config.Ctx, key, cartTTL).Err()
}

func setCartItems(key string, items []dto.RequestCartItems) error {
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	err = config.RedisClient.Set(config.Ctx, key, data, cartTTL).Err()
	if err != nil {
		log.Printf("Redis SetCart error: %v", err)
	}
	return err
}
