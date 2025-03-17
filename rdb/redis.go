package rdb

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"

	"github.com/virzz/mulan/utils"
	"github.com/virzz/vlog"
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
		vlog.Info(cmd.String())
		next(ctx, cmd)
		return nil
	}
}

// 执行管道命令时调用
func (DebugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

var (
	rdb      *redis.Client
	oncePlus utils.OncePlus
	Nil      = redis.Nil
)

func connect(cfg *Config) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Pass,
		DB:       cfg.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			vlog.Info("Redis is connected")
			return nil
		},
	})
	if cfg.Debug {
		rdb.AddHook(DebugHook{})
	}
	return rdb.Ping(context.Background()).Err()
}

func Init(cfg *Config, force ...bool) error {
	if len(force) > 0 && force[0] {
		return connect(cfg)
	}
	return oncePlus.Do(func() (err error) {
		return connect(cfg)
	})
}

func R() *redis.Client { return rdb }
