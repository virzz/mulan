package app_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

func TestApp(t *testing.T) {
	Conf = &app.Config{}
	meta := &app.Meta{
		ID:          "com.virzz.mulan.example",
		Name:        "example",
		Description: "ExampleService",
		Version:     Version,
		Commit:      Commit,
	}
	std := app.New(meta)

	web.SetVersionHandler(meta.Name, meta.Version, meta.Commit)

	applyFunc := func(api *gin.RouterGroup) {
		api.Handle("GET", "/", func(c *gin.Context) {
			c.String(200, "Hello, World!")
		})
	}

	std.SetPreInit(func(ctx context.Context) error { return nil })
	std.SetValidate(func() error { return nil })

	std.SetAction(func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		httpCfg := Conf.GetHTTP().WithRequestID(true)
		httpSrv, err := web.New(httpCfg, applyFunc)
		if err != nil {
			return err
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		go func() {
			err := httpSrv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Error("Failed to run http server", zap.Error(err))
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
			httpSrv.Close()
		case syscall.SIGTERM:
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			httpSrv.Shutdown(ctx)
		}
		return nil
	})
	if err := std.Execute(context.Background(), Conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
