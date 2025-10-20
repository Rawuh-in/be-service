package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) *Redis {
	fmt.Println("Redis connecting to: ", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

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
