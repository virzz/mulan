package app

import (
	"context"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/virzz/mulan/db"
	"github.com/virzz/mulan/log"
	"github.com/virzz/mulan/web"
)

type (
	ActionFunc   func(cmd *cobra.Command, args []string) error
	PreInitFunc  func(context.Context) error
	ValidateFunc func() error

	Meta struct {
		ID          string
		Name        string
		Description string
		Version     string
		Commit      string
		BuildAt     string
	}
	App struct {
		*Meta
		rootCmd     *cobra.Command
		action      ActionFunc
		routers     *web.Routers
		preInit     PreInitFunc
		validate    ValidateFunc
		conf        Configer
		log         *zap.Logger
		remote      *Remote //lint:ignore U1000 remote config
		maintainCmd map[string]*db.Config
	}
)

var (
	std  *App
	Conf Configer
)

func New(meta *Meta) *App {
	std = &App{
		Meta:        meta,
		conf:        &Config{},
		validate:    nil,
		maintainCmd: map[string]*db.Config{},
		rootCmd: &cobra.Command{
			CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
			SilenceErrors:     true,
			SilenceUsage:      true,
			RunE: func(_ *cobra.Command, _ []string) error {
				panic("execute action not implemented")
			},
		},
	}
	return std
}

func (app *App) AddMaintainCmd(name string, cfg *db.Config) {
	if app.maintainCmd == nil {
		app.maintainCmd = map[string]*db.Config{}
	}
	app.maintainCmd[name] = cfg
}
func (app *App) DisableConfigCmd()                           { disableConfigCmd = true }
func (app *App) SetPreInit(f PreInitFunc)                    { app.preInit = f }
func (app *App) SetValidate(f ValidateFunc)                  { app.validate = f }
func (app *App) SetRouters(routers *web.Routers)             { app.routers = routers }
func (app *App) SetConfig(config Configer)                   { app.conf = config }
func (app *App) Register(f web.RegisterFunc)                 { app.routers.Register(f) }
func (app *App) AddCommand(cmd ...*cobra.Command)            { app.rootCmd.AddCommand(cmd...) }
func (app *App) RootCmd() *cobra.Command                     { return app.rootCmd }
func (app *App) Conf() Configer                              { return app.conf }
func (app *App) Routers() *web.Routers                       { return app.routers }
func (app *App) Run(ctx context.Context, cfg Configer) error { return app.Execute(ctx, cfg) }

func (app *App) AddFlagSet(fs ...*pflag.FlagSet) {
	for _, f := range fs {
		app.rootCmd.Flags().AddFlagSet(f)
	}
}
func (app *App) preRunE() (err error) {
	if app.conf != nil {
		err = viper.Unmarshal(app.conf, func(dc *mapstructure.DecoderConfig) { dc.TagName = "json" })
		if err != nil {
			return err
		}
	}
	logger, err := log.NewWithConfig(app.conf.GetLog())
	if err != nil {
		return err
	}
	app.log = logger.Named("app")
	return nil
}

func (app *App) SetAction(action ActionFunc) {
	app.action = func(cmd *cobra.Command, args []string) error {
		if app.preInit != nil {
			err := app.preInit(cmd.Context())
			if err != nil {
				return err
			}
		}
		return action(cmd, args)
	}
}
func (app *App) Execute(ctx context.Context, cfg Configer) error {
	// Config
	app.rootCmd.PersistentFlags().CountP("verbose", "v", "verbose mode")
	app.rootCmd.PersistentFlags().String("instance", "default", "instance name")
	app.rootCmd.PersistentFlags().String("config", "", "config file")
	app.rootCmd.PersistentFlags().AddFlagSet(log.FlagSet())
	logger, err := log.New(zapcore.DPanicLevel, true, app.Name)
	if err != nil {
		return err
	}
	defer logger.Sync()
	app.log = logger.Named("app")
	// Action
	if app.action != nil {
		app.rootCmd.RunE = app.action
	}
	err = viper.BindPFlags(app.rootCmd.PersistentFlags())
	if err != nil {
		return err
	}
	err = viper.BindPFlags(app.rootCmd.Flags())
	if err != nil {
		return err
	}

	viper.SetEnvPrefix(app.Name)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	app.conf = cfg
	return app.ExecuteE(ctx)
}
