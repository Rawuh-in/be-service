package redis

import "github.com/redis/go-redis/v9"

type RedisManager struct {
	Client   *redis.Client
	Prefixes *KeyPrefixes
}

func (r *RedisManager) SetClient(client *redis.Client) {
	r.Client = client
}

func (r *RedisManager) SetPrefixes(p *KeyPrefixes) {
	r.Prefixes = p
}
