package app

import (
	"context"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	ActionFunc   func(cmd *cobra.Command, args []string) error
	PreInitFunc  func(context.Context) error
	ValidateFunc func() error

	Meta struct {
		ID          string
		Name        string
		CNName      string
		Description string
		Version     string
		Commit      string
		BuildAt     string
	}
	App struct {
		*Meta
		rootCmd  *cobra.Command
		action   ActionFunc
		preInit  PreInitFunc
		validate ValidateFunc
		conf     Configer
		log      *zap.Logger
		remote   *Remote //lint:ignore U1000 remote config
	}
)

var (
	std  *App
	Conf Configer
)

func New(meta *Meta) *App {
	std = &App{
		Meta:     meta,
		conf:     &Config{},
		validate: nil,
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

func (app *App) DisableConfigCmd()                           { disableConfigCmd = true }
func (app *App) SetPreInit(f PreInitFunc)                    { app.preInit = f }
func (app *App) SetValidate(f ValidateFunc)                  { app.validate = f }
func (app *App) SetConfig(config Configer)                   { app.conf = config }
func (app *App) SetLogger(log *zap.Logger)                   { app.log = log }
func (app *App) AddCommand(cmd ...*cobra.Command)            { app.rootCmd.AddCommand(cmd...) }
func (app *App) RootCmd() *cobra.Command                     { return app.rootCmd }
func (app *App) Conf() Configer                              { return app.conf }
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

func (app *App) Execute(ctx context.Context, cfg Configer) (err error) {
	// Config
	app.rootCmd.PersistentFlags().CountP("verbose", "v", "verbose mode")
	app.rootCmd.PersistentFlags().String("instance", "default", "instance name")
	app.rootCmd.PersistentFlags().String("config", "", "config file")

	app.log, err = zap.NewProduction()
	if err != nil {
		return err
	}
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
