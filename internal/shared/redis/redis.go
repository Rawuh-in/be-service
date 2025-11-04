package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) *Redis {
	fmt.Println("Redis connecting to:", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{client}
}

func NewRedisFromDSN(dsn string) *Redis {
	if dsn == "" {
		fmt.Println("No REDIS_DSN provided, falling back to localhost:6379")
		return NewRedis("localhost:6379", "", 0)
	}

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		fmt.Println("Failed to parse REDIS_DSN:", err)
		return NewRedis("localhost:6379", "", 0)
	}

	// ✅ Enable TLS automatically if rediss:// scheme
	if strings.HasPrefix(dsn, "rediss://") {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // skip cert verification (solves Heroku hostname mismatch)
		}
		fmt.Println("Redis TLS enabled (InsecureSkipVerify=true)")
	}

	client := redis.NewClient(opts)
	fmt.Println("Redis connecting via DSN:", opts.Addr)
	return &Redis{client}
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", redis.Nil
		}
		return "", fmt.Errorf("failed to get key: %s", err)
	}
	return val, nil
}

func (r *Redis) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %s", err)
	}
	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key: %s", err)
	}
	return nil
}

func (r *Redis) Update(ctx context.Context, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %s", err)
	}
	if err := r.client.Set(ctx, key, data, redis.KeepTTL).Err(); err != nil {
		return fmt.Errorf("failed to update key: %s", err)
	}
	return nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	if _, err := r.client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to delete key: %s", err)
	}
	return nil
}
