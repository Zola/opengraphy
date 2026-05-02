package stats

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func AllowFixedWindow(ctx context.Context, rdb *redis.Client, key string, limit int, window time.Duration) bool {
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return true
	}
	if count == 1 {
		_ = rdb.Expire(ctx, key, window).Err()
	}
	return count <= int64(limit)
}
