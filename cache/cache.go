package cache

import (
	"context"
	"time"
)

// Cache describes the minimal cache operations used by the database service.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}
