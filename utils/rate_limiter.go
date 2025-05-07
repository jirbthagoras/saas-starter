package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

var ctx = context.Background()

func AllowRequest(rdb *redis.Client, apiKey string, maxRequest int, window time.Duration) (bool, int64, error) {
	key := "rate_limit:" + apiKey
	now := time.Now().Unix()

	// Define the windows
	windowStart := now - int64(window.Seconds())

	// Starts a pipeline
	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// Remove old timestamps
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

		// Creates a new timestamp
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: now,
		})

		// Set the Expire to auto cleanup
		remainingTTL := window - time.Duration(now-windowStart)*time.Second
		if remainingTTL > 0 {
			pipe.Expire(ctx, key, remainingTTL)
		}
		return nil
	})

	if err != nil {
		return false, 0, err
	}

	slog.Debug("Created a new timestamp", "apiKey", apiKey, "timestamp", now)

	// Checking the Sorted Sets
	count, err := rdb.ZCount(ctx, key, "-inf", "+inf").Result()

	if err != nil {
		return false, 0, err
	}

	// Returns the results
	return count <= int64(maxRequest), count, nil
}
