package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"rijig/config"

	"github.com/go-redis/redis/v8"
)

func SetCache(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	err = config.RedisClient.Set(config.Ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

func GetCache(key string, dest interface{}) error {
	val, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("ErrCacheMiss")
		}
		return fmt.Errorf("failed to get cache: %v", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %v", err)
	}

	return nil
}

func DeleteCache(key string) error {
	err := config.RedisClient.Del(config.Ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %v", err)
	}
	return nil
}

func ScanAndDelete(pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := config.RedisClient.Scan(config.Ctx, cursor, pattern, 10).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := config.RedisClient.Del(config.Ctx, keys...).Err(); err != nil {
				return err
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}

func CacheExists(key string) (bool, error) {
	exists, err := config.RedisClient.Exists(config.Ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %v", err)
	}
	return exists > 0, nil
}

func SetCacheWithTTL(key string, value interface{}, expiration time.Duration) error {
	return SetCache(key, value, expiration)
}

func GetTTL(key string) (time.Duration, error) {
	ttl, err := config.RedisClient.TTL(config.Ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %v", err)
	}
	return ttl, nil
}

func RefreshTTL(key string, expiration time.Duration) error {
	err := config.RedisClient.Expire(config.Ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh TTL: %v", err)
	}
	return nil
}

func IncrementCounter(key string, expiration time.Duration) (int64, error) {
	val, err := config.RedisClient.Incr(config.Ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter: %v", err)
	}

	if val == 1 {
		config.RedisClient.Expire(config.Ctx, key, expiration)
	}

	return val, nil
}

func DecrementCounter(key string) (int64, error) {
	val, err := config.RedisClient.Decr(config.Ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement counter: %v", err)
	}
	return val, nil
}

func GetCounter(key string) (int64, error) {
	val, err := config.RedisClient.Get(config.Ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get counter: %v", err)
	}
	return val, nil
}

func SetList(key string, values []interface{}, expiration time.Duration) error {

	config.RedisClient.Del(config.Ctx, key)

	if len(values) > 0 {
		err := config.RedisClient.LPush(config.Ctx, key, values...).Err()
		if err != nil {
			return fmt.Errorf("failed to set list: %v", err)
		}

		if expiration > 0 {
			config.RedisClient.Expire(config.Ctx, key, expiration)
		}
	}

	return nil
}

func GetList(key string) ([]string, error) {
	vals, err := config.RedisClient.LRange(config.Ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %v", err)
	}
	return vals, nil
}

func AddToList(key string, value interface{}, expiration time.Duration) error {
	err := config.RedisClient.LPush(config.Ctx, key, value).Err()
	if err != nil {
		return fmt.Errorf("failed to add to list: %v", err)
	}

	if expiration > 0 {
		config.RedisClient.Expire(config.Ctx, key, expiration)
	}

	return nil
}

func SetHash(key string, fields map[string]interface{}, expiration time.Duration) error {
	err := config.RedisClient.HMSet(config.Ctx, key, fields).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash: %v", err)
	}

	if expiration > 0 {
		config.RedisClient.Expire(config.Ctx, key, expiration)
	}

	return nil
}

func GetHash(key string) (map[string]string, error) {
	vals, err := config.RedisClient.HGetAll(config.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash: %v", err)
	}
	return vals, nil
}

func GetHashField(key, field string) (string, error) {
	val, err := config.RedisClient.HGet(config.Ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("hash field not found")
		}
		return "", fmt.Errorf("failed to get hash field: %v", err)
	}
	return val, nil
}

func SetHashField(key, field string, value interface{}, expiration time.Duration) error {
	err := config.RedisClient.HSet(config.Ctx, key, field, value).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash field: %v", err)
	}

	if expiration > 0 {
		config.RedisClient.Expire(config.Ctx, key, expiration)
	}

	return nil
}

func FlushDB() error {
	err := config.RedisClient.FlushDB(config.Ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush database: %v", err)
	}
	return nil
}

func GetAllKeys(pattern string) ([]string, error) {
	keys, err := config.RedisClient.Keys(config.Ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %v", err)
	}
	return keys, nil
}
