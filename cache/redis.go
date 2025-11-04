package cache

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache provides a simple Cache implementation backed by Redis.
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a Redis-backed cache with an optional key prefix.
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

func (c *RedisCache) namespaced(key string) string {
	if c.prefix == "" {
		return key
	}
	if strings.HasSuffix(c.prefix, ":") {
		return c.prefix + key
	}
	return c.prefix + ":" + key
}

// Get retrieves a value from Redis. A missing key is returned as an empty string with nil error.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	res, err := c.client.Get(ctx, c.namespaced(key)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

// Set stores a value in Redis for the given TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, c.namespaced(key), value, ttl).Err()
}

// Delete removes one or more keys from Redis.
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if len(keys) == 1 {
		return c.client.Del(ctx, c.namespaced(keys[0])).Err()
	}
	namespaced := make([]string, len(keys))
	for i, k := range keys {
		namespaced[i] = c.namespaced(k)
	}
	return c.client.Del(ctx, namespaced...).Err()
}
