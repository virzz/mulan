package app

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	slogzap "github.com/samber/slog-zap"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/virzz/mulan/service"
)

type (
	ActionFunc  func(cmd *cobra.Command, args []string) error
	PreInitFunc func(context.Context) error

	Meta struct {
		ID          string
		Name        string
		CNName      string
		Description string
		Version     string
		Commit      string
		BuildAt     string
	}
	App[T any] struct {
		*Meta
		rootCmd *cobra.Command
		preInit PreInitFunc
		log     *zap.Logger
		remote  *Remote
		conf    *T
		srvs    []service.Servicer
	}
)

var replacer = strings.NewReplacer(".", "__", "-", "__")

func New[T any](meta *Meta, cfg *T) *App[T] {
	_ = godotenv.Load()
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	std := &App[T]{
		Meta: meta,
		conf: cfg,
		log:  zap.L(),
		srvs: make([]service.Servicer, 0),
		rootCmd: &cobra.Command{
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(_ *cobra.Command, _ []string) error {
				panic("execute action not implemented")
			},
		},
	}
	std.internalCmd()
	return std
}

func (app *App[T]) SetPreInit(f PreInitFunc) *App[T]  { app.preInit = f; return app }
func (app *App[T]) SetLogger(log *zap.Logger) *App[T] { app.log = log; return app }
func (app *App[T]) RootCmd() *cobra.Command           { return app.rootCmd }
func (app *App[T]) Conf() *T                          { return app.conf }

func (app *App[T]) AddFlagSet(fs ...*pflag.FlagSet) *App[T] {
	flags := app.rootCmd.Flags()
	for _, f := range fs {
		flags.AddFlagSet(f)
	}
	return app
}

func (app *App[T]) AddService(srvs ...service.Servicer) *App[T] {
	app.srvs = append(app.srvs, srvs...)
	return app
}

func (app *App[T]) AddCommand(cmd ...*cobra.Command) *App[T] {
	internalCmds := []string{"config", "version"}
	for _, c := range cmd {
		if slices.Contains(internalCmds, c.Use) {
			app.log.Warn("internal command cannot be added", zap.String("command", c.Use))
			continue
		}
		app.rootCmd.AddCommand(c)
	}
	return app
}

func loadConfig[T any](app *App[T]) func(cmd *cobra.Command, args []string) (err error) {
	logger := app.log.Named("viper")
	return func(cmd *cobra.Command, args []string) (err error) {
		fs := cmd.Flags()
		instance, _ := fs.GetString("instance")
		configPath, _ := fs.GetString("config")
		viper.SetOptions(
			viper.WithLogger(
				slog.New(
					slogzap.Option{
						Level:  slog.LevelWarn,
						Logger: logger,
					}.NewZapHandler(),
				),
			),
		)
		configLoaded := false
		configSource := "env"
		if configPath != "" {
			configSource = "args"
			viper.SetConfigFile(configPath)
		} else if app.remote != nil {
			app.remote.instance = instance
			if err = viper.ReadRemoteConfig(); err == nil {
				configSource = "remote"
				configLoaded = true
			}
		}
		if !configLoaded {
			viper.AddConfigPath(".")
			viper.SetConfigName("config_" + instance)
			viper.SetConfigType("")
			if viper.ReadInConfig() == nil {
				configSource = viper.ConfigFileUsed()
			}
		}
		// Viper config unmarshal to app.conf
		err = viper.Unmarshal(app.conf, func(dc *mapstructure.DecoderConfig) { dc.TagName = "json" })
		if err != nil {
			logger.Error("Failed to unmarshal config", zap.Error(err))
			return err
		}
		logger.Info("Config loaded from", zap.String("source", configSource))
		return nil
	}
}

func (app *App[T]) Execute(ctx context.Context, action ...ActionFunc) (err error) {
	// Config
	app.rootCmd.PersistentFlags().StringP("instance", "i", "default", "instance name")
	app.rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	err = viper.BindPFlags(app.rootCmd.PersistentFlags())
	if err != nil {
		return err
	}
	err = viper.BindPFlags(app.rootCmd.Flags())
	if err != nil {
		return err
	}
	viper.SetEnvPrefix(app.Name)
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	app.rootCmd.PreRunE = loadConfig(app)

	// Action
	var currentAction ActionFunc
	if len(action) > 0 && action[0] != nil {
		currentAction = action[0]
	} else {
		currentAction = defaultAction(app)
	}
	app.rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if app.preInit != nil {
			err := app.preInit(cmd.Context())
			if err != nil {
				return err
			}
		}
		return currentAction(cmd, args)
	}
	return app.rootCmd.ExecuteContext(ctx)
}
