package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/gabrzb/auth-go-gin/internal/config"
	"github.com/redis/go-redis/v9"
)

const keyPrefix = "blk:"

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(cfg *config.Config) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connected")
	return &RedisStore{client: client}, nil
}

func (r *RedisStore) Add(token string, ttl time.Duration) error {
	return r.client.Set(context.Background(), hashToken(token), 1, ttl).Err()
}

func (r *RedisStore) Contains(ctx context.Context, token string) (bool, error) {
	n, err := r.client.Exists(ctx, hashToken(token)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// hashToken returns a non-sensitive Redis key so raw JWTs are never stored.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return keyPrefix + hex.EncodeToString(h[:])
}
