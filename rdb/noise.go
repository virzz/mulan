package rdb

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func NoiseSet(ctx context.Context, key string, expire time.Duration) int64 {
	if key == "" {
		return 0
	}
	_, err := rdb.SetEx(ctx, key, 1, expire).Result()
	if err != nil {
		zap.L().Error("Failed to incr noise", zap.String("key", key), zap.Error(err))
		return 0
	}
	return 1
}

func NoiseGet(ctx context.Context, key string) int64 {
	if key == "" {
		return 0
	}
	n, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		zap.L().Error("Failed to incr noise", zap.String("key", key), zap.Error(err))
	}
	return n
}

func NoiseCheck(ctx context.Context, key string, count int64, expire time.Duration) (bool, int64) {
	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		zap.L().Warn("Failed to get noise", zap.String("key", key), zap.Error(err))
	}
	if n == 0 {
		_, err = rdb.SetEx(ctx, key, 1, expire).Result()
		if err != nil {
			zap.L().Error("Failed to set noise", zap.String("key", key), zap.Error(err))
		}
		return true, 1
	}
	n, err = rdb.Incr(ctx, key).Result()
	if err != nil {
		zap.L().Error("Failed to incr noise", zap.String("key", key), zap.Error(err))
	}
	if n > count {
		rdb.Del(ctx, key)
		return false, n
	}
	return true, n
}
