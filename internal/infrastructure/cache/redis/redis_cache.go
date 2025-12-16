package redis

import (
	"context"
	"time"

	"soccer-manager-api/internal/ports/cache"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}


func NewRedisCache(client *redis.Client) cache.Cache {
	return &redisCache{client: client}
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return []byte(val), nil
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, ttl int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

