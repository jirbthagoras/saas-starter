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
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		// Remove old timestamps
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

		// Creates a new timestamp
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: now,
		})
		slog.Debug("Created a new timestamp")

		// Set the Expire to auto cleanup
		pipe.Expire(ctx, key, window)
		return nil
	})

	if err != nil {
		return false, 0, err
	}

	// Checking the Sorted Sets
	count, err := rdb.ZCount(ctx, key, fmt.Sprintf("%d", windowStart), fmt.Sprintf("%d", now)).Result()

	if err != nil {
		return false, 0, err
	}

	// Returns the results
	return count <= int64(maxRequest), count, nil
}
