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
	"github.com/virzz/mulan/db"
	"github.com/virzz/mulan/rdb"
	"github.com/virzz/mulan/web"
	"go.uber.org/zap"
)

type Config struct {
	HTTP  web.Config `json:"http" yaml:"http"`
	Token web.Token  `json:"token" yaml:"token"`
	DB    db.Config  `json:"db" yaml:"db"`
	RDB   rdb.Config `json:"rdb" yaml:"rdb"`
}

var (
	Version string = "1.0.0"
	Commit  string = "dev"

	Conf *Config
)

func Example() {
	meta := &app.Meta{
		ID:          "com.virzz.mulan.example",
		Name:        "example",
		Description: "ExampleService",
		Version:     Version,
		Commit:      Commit,
	}
	std := app.New(meta, nil)
	std.SetPreInit(func(ctx context.Context) error {
		return nil
	})
	std.SetValidate(func() error {
		return nil
	})
	web.SetVersionHandler(meta.Name, meta.Version, meta.Commit)

	applyFunc := func(api *gin.RouterGroup) {
		api.Handle("GET", "/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello, World!"})
		})
	}

	std.SetAction(func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		httpCfg := Conf.HTTP.WithRequestID(true)
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
