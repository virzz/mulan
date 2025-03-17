package rdb

import (
	"context"
	"time"

	"github.com/virzz/vlog"
)

func NoiseSet(ctx context.Context, key string, expire time.Duration) int64 {
	if key == "" {
		return 0
	}
	_, err := rdb.SetEx(ctx, key, 1, expire).Result()
	if err != nil {
		vlog.Error("Failed to incr noise", "key", key, "err", err.Error())
	}
	return 1
}

func NoiseGet(ctx context.Context, key string) int64 {
	if key == "" {
		return 0
	}
	n, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		vlog.Error("Failed to incr noise", "key", key, "err", err.Error())
	}
	return n
}

func NoiseCheck(ctx context.Context, key string, count int64, expire time.Duration) (bool, int64) {
	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		vlog.Warn("Failed to get noise", "key", key, "err", err.Error())
	}
	if n == 0 {
		_, err = rdb.SetEx(ctx, key, 1, expire).Result()
		if err != nil {
			vlog.Error("Failed to set noise", "key", key, "err", err.Error())
		}
		return true, 1
	}
	n, err = rdb.Incr(ctx, key).Result()
	if err != nil {
		vlog.Error("Failed to incr noise", "key", key, "err", err.Error())
	}
	if n > count {
		rdb.Del(ctx, key)
		return false, n
	}
	return true, n
}
