package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) *Redis {
	fmt.Println("Redis connecting to:", addr)
	client := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  password,
		DB:        db,
		TLSConfig: &tls.Config{}, // ✅ enable TLS manually for rediss-like endpoints
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
		fmt.Println("failed to parse REDIS_DSN, falling back to addr form:", err)
		return NewRedis("localhost:6379", "", 0)
	}

	// ✅ force TLS when using rediss://
	if opts.TLSConfig == nil {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			// You can disable cert verification if needed (not recommended for production):
			// In production, better to use proper CA cert.
			InsecureSkipVerify: true,
		}
	}

	fmt.Println("Redis connecting via DSN with TLS")
	client := redis.NewClient(opts)
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

	err = r.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %s", err)
	}

	return nil
}

func (r *Redis) Update(ctx context.Context, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %s", err)
	}

	err = r.client.Set(ctx, key, data, redis.KeepTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %s", err)
	}

	return nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to del key: %s", err)
	}

	return nil
}
