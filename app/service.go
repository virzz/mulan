package app

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (app *App[T]) Serve() error {
	wg := errgroup.Group{}
	for _, srv := range app.srvs {
		wg.Go(srv.Serve)
	}
	return wg.Wait()
}

func (app *App[T]) Close() error {
	wg := errgroup.Group{}
	for _, srv := range app.srvs {
		wg.Go(srv.Close)
	}
	return wg.Wait()
}

func (app *App[T]) Shutdown(ctx context.Context) error {
	wg := errgroup.Group{}
	for _, srv := range app.srvs {
		wg.Go(func() error { return srv.Shutdown(ctx) })
	}
	return wg.Wait()
}
