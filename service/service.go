package service

import (
	"context"
)

type Servicer interface {
	Serve() error
	Shutdown(ctx context.Context) error
	Close() error
}
