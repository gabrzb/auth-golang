package store

import (
	"context"
	"log"
	"time"

	"github.com/gabrzb/auth-go-gin/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(cfg *config.Config) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected")
	return &RedisStore{client: client}
}

func (r *RedisStore) Add(token string, ttl time.Duration) error {
	return r.client.Set(context.Background(), token, 1, ttl).Err()
}

func (r *RedisStore) Contains(token string) (bool, error) {
	n, err := r.client.Exists(context.Background(), token).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
