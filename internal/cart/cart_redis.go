package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rijig/config"

	"github.com/go-redis/redis/v8"
)

const CartTTL = 30 * time.Minute
const CartKeyPrefix = "cart:"

func buildCartKey(userID string) string {
	return fmt.Sprintf("%s%s", CartKeyPrefix, userID)
}

func SetCartToRedis(ctx context.Context, userID string, cart RequestCartDTO) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, buildCartKey(userID), data, CartTTL).Err()
}

func RefreshCartTTL(ctx context.Context, userID string) error {
	return config.RedisClient.Expire(ctx, buildCartKey(userID), CartTTL).Err()
}

func GetCartFromRedis(ctx context.Context, userID string) (*RequestCartDTO, error) {
	val, err := config.RedisClient.Get(ctx, buildCartKey(userID)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var cart RequestCartDTO
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		return nil, err
	}
	return &cart, nil
}

func DeleteCartFromRedis(ctx context.Context, userID string) error {
	return config.RedisClient.Del(ctx, buildCartKey(userID)).Err()
}

func GetExpiringCartKeys(ctx context.Context, threshold time.Duration) ([]string, error) {
	keys, err := config.RedisClient.Keys(ctx, CartKeyPrefix+"*").Result()
	if err != nil {
		return nil, err
	}

	var expiringKeys []string
	for _, key := range keys {
		ttl, err := config.RedisClient.TTL(ctx, key).Result()
		if err != nil {
			continue
		}

		if ttl > 0 && ttl <= threshold {
			expiringKeys = append(expiringKeys, key)
		}
	}

	return expiringKeys, nil
}
