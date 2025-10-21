package app_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/virzz/mulan/app"
	"github.com/virzz/mulan/web"
	"go.uber.org/zap"
)

func testAction[T any](app *app.App[T]) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		go func() {
			err := app.Serve()
			if err != nil {
				zap.L().Error("Failed to run server", zap.Error(err))
				sig <- os.Interrupt
			}
		}()

		// close server after 10 seconds for test
		go func() {
			<-time.After(10 * time.Second)
			sig <- os.Interrupt
		}()

		switch <-sig {
		case os.Interrupt:
			app.Close()
		case syscall.SIGTERM:
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			app.Shutdown(ctx)
		}
		return nil
	}
}

func TestApp(t *testing.T) {
	Conf = &Config{}
	meta := &app.Meta{
		ID:          "com.virzz.mulan.example",
		Name:        "example",
		Description: "ExampleService",
		Version:     Version,
		Commit:      Commit,
		BuildAt:     BuildAt,
	}
	std := app.New(meta, Conf)

	webSrv := web.New(
		&Conf.HTTP,
		&web.Info{Name: meta.Name, Version: meta.Version, Commit: meta.Commit},
		func(api gin.IRouter) {
			api.Handle("GET", "/", func(c *gin.Context) {
				c.String(200, "Hello, World!")
			})
		},
	)

	std.AddService(webSrv)

	if err := std.Execute(context.Background(), testAction(std)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
