package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rijig/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user *model.User) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveUser(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

type RedisRepository interface {
	StoreData(key string, data interface{}, expiration time.Duration) error
	GetData(key string) (interface{}, error) // Mengembalikan interface{}
	DeleteData(key string) error
}

type redisRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisRepository membuat instance baru dari redisRepository
func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

// StoreData menyimpan data ke dalam Redis (dalam format JSON)
func (r *redisRepository) StoreData(key string, data interface{}, expiration time.Duration) error {
	// Marshaling data ke JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Simpan JSON ke Redis
	err = r.client.Set(r.ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store data in Redis: %v", err)
	}
	return nil
}

// GetData mengambil data dari Redis berdasarkan key
func (r *redisRepository) GetData(key string) (interface{}, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get data from Redis: %v", err)
	}

	// Unmarshal data JSON kembali ke objek
	var data interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return data, nil
}

// DeleteData menghapus data di Redis berdasarkan key
func (r *redisRepository) DeleteData(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete data from Redis: %v", err)
	}
	return nil
}
