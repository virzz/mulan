package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/virzz/daemon/v2"
	"github.com/virzz/vlog"

	"github.com/virzz/mulan/db"
	"github.com/virzz/mulan/web"
)

type (
	PreInitFunc  func(context.Context) error
	ValidateFunc func() error
	App          struct{ ID, Name, Description, Version, Commit string }
)

var (
	std      *App
	routers  web.Routers
	preInit  PreInitFunc
	validate ValidateFunc
	Conf     Configer

	cmds   []*cobra.Command
	models []any
)

func (app *App) Run(ctx context.Context, cfg Configer) error {
	std = app
	return Execute(context.Background(), cfg)
}

func Execute(ctx context.Context, cfg Configer) error {
	daemon.New(std.ID, std.Name, std.Description, std.Version, std.Commit)
	daemon.RegisterConfig(cfg)
	daemon.AddCommand(&cobra.Command{
		Use: "validate", Aliases: []string{"valid"},
		Short: "Validate Configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if validate == nil {
				return errors.New("not implemented validate function")
			}
			if err := cfg.Validate(); err != nil {
				return err
			}
			if err := validate(); err != nil {
				return err
			}
			fmt.Println("Valid configuration")
			return nil
		},
	})

	daemon.RootCmd().AddGroup(&cobra.Group{ID: "maintain", Title: "Maintain Commands"})
	// Maintain Cmds
	daemon.AddCommand(db.MaintainCommand(cfg.GetDB())...)
	daemon.AddCommand(cmds...)
	daemon.Execute(func(cmd *cobra.Command, _ []string) error {
		os.MkdirAll("logs", 0755)
		vlog.New(filepath.Join("logs", std.Name+".log"))
		vlog.Log = vlog.Log.With("service", std.Name)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if preInit != nil {
			if err := preInit(ctx); err != nil {
				return err
			}
		}
		webCfg := cfg.GetHTTP().Check().
			WithVersion(std.Version).WithCommit(std.Commit)
		httpSrv, err := web.New(webCfg, routers, []gin.HandlerFunc{web.LogMw}, nil)
		if err != nil {
			return err
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		go func() {
			err := httpSrv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Fprintln(os.Stderr, err)
				sig <- os.Interrupt
			}
		}()

		switch <-sig {
		case os.Interrupt:
			httpSrv.Close()
		case syscall.SIGTERM:
			httpSrv.Shutdown(ctx)
		}
		return nil
	})
	return nil
}

func ID() string                       { return std.ID }
func AppID() string                    { return std.ID }
func Name() string                     { return std.Name }
func Description() string              { return std.Description }
func Desc() string                     { return std.Description }
func Version() string                  { return std.Version }
func Commit() string                   { return std.Commit }
func SetPreInit(f PreInitFunc)         { preInit = f }
func SetValidate(f ValidateFunc)       { validate = f }
func Register(f web.RegisterFunc)      { routers.Register(f) }
func RegisterModels(ms ...any)         { models = ms }
func AddCommand(cmd ...*cobra.Command) { cmds = append(cmds, cmd...) }

func Run(ctx context.Context, app *App, cfg Configer) error {
	return std.Run(ctx, cfg)
}
