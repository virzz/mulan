package app_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/virzz/mulan/app"
	"github.com/virzz/mulan/web"
	"go.uber.org/zap"
)

type Config struct {
	//lint:ignore SA5008 Ignore JSON option "squash"
	app.Config `json:",inline,squash" yaml:",inline"`
}

var (
	Version string = "1.0.0"
	Commit  string = "dev"

	Conf app.Configer
)

func Example() {
	meta := &app.Meta{
		ID:          "com.virzz.mulan.example",
		Name:        "example",
		Description: "ExampleService",
		Version:     Version,
		Commit:      Commit,
	}
	std := app.New(meta)
	std.SetPreInit(func(ctx context.Context) error {
		return nil
	})
	std.SetValidate(func() error {
		return nil
	})
	web.SetVersionHandler(meta.Name, meta.Version, meta.Commit)
	routers := web.NewRouters()
	routers.Handle("GET", "/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	std.SetRouters(routers)
	std.SetAction(func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		httpCfg := Conf.GetHTTP().WithRequestID(true)
		httpSrv, err := web.New(httpCfg, std.Routers(), nil, nil)
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
		panic(err)
	}
}
