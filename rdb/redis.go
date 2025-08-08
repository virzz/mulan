package rdb

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type DebugHook struct{}

// 当创建网络连接时调用
func (DebugHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// 执行命令时调用
func (DebugHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		zap.L().Info(cmd.String())
		return next(ctx, cmd)
	}
}

// 执行管道命令时调用
func (DebugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

var (
	rdb *redis.Client
	Nil = redis.Nil
)

func New(cfg *Config) (*redis.Client, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Pass,
		DB:       cfg.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			zap.L().Info("Redis is connected")
			return nil
		},
	})
	if cfg.Debug {
		rdb.AddHook(DebugHook{})
	}
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
