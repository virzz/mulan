package app

import (
	"context"
	"log/slog"
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
	App struct {
		*Meta
		debug   int
		rootCmd *cobra.Command
		preInit PreInitFunc
		log     *zap.Logger
		remote  *Remote
		conf    any
		srvs    []service.Servicer
	}
)

var replacer = strings.NewReplacer(".", "__", "-", "__")

func New(meta *Meta, cfg any) *App {
	_ = godotenv.Load()
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	std := &App{
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

func (app *App) SetPreInit(f PreInitFunc) *App  { app.preInit = f; return app }
func (app *App) SetLogger(log *zap.Logger) *App { app.log = log; return app }
func (app *App) RootCmd() *cobra.Command        { return app.rootCmd }
func (app *App) Conf() any                      { return app.conf }

func (app *App) AddFlagSet(fs ...*pflag.FlagSet) *App {
	flags := app.rootCmd.Flags()
	for _, f := range fs {
		flags.AddFlagSet(f)
	}
	return app
}

func (app *App) AddService(srvs ...service.Servicer) *App {
	app.srvs = append(app.srvs, srvs...)
	return app
}

func (app *App) AddCommand(cmd ...*cobra.Command) *App {
	for _, c := range cmd {
		if c.Use == "version" {
			app.log.Warn("internal command cannot be added", zap.String("command", c.Use))
			continue
		}
		app.rootCmd.AddCommand(c)
	}
	return app
}

func loadConfig(app *App) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		fs := cmd.Flags()
		configPath, _ := fs.GetString("config")
		viper.SetOptions(viper.WithLogger(slog.New(
			slogzap.Option{Level: slog.LevelWarn, Logger: app.log}.NewZapHandler(),
		)))
		configLoaded := false
		configSource := "env"
		if configPath != "" {
			configSource = "args"
			viper.SetConfigFile(configPath)
		} else if app.remote != nil {
			if err = viper.ReadRemoteConfig(); err == nil {
				configSource = "remote"
				configLoaded = true
			}
		}
		if viper.ConfigFileUsed() == "" {
			viper.AddConfigPath(".")
			viper.SetConfigName("config")
			viper.SetConfigType("")
		}
		if !configLoaded {
			if viper.ReadInConfig() == nil {
				configSource = viper.ConfigFileUsed()
			}
		}
		// Viper config unmarshal to app.conf
		err = viper.Unmarshal(app.conf, func(dc *mapstructure.DecoderConfig) { dc.TagName = "json" })
		if err != nil {
			app.log.Error("Failed to unmarshal config", zap.Error(err))
			return err
		}
		app.log.Info("Config loaded from", zap.String("source", configSource))
		return nil
	}
}

func (app *App) Execute(ctx context.Context, action ...ActionFunc) (err error) {
	// Config
	app.rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	app.rootCmd.PersistentFlags().CountVarP(&app.debug, "debug", "d", "debug mode")
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
	app.log.Info("App Info",
		zap.String("id", app.ID),
		zap.String("name", app.Name),
		zap.String("cn_name", app.CNName),
		zap.String("description", app.Description),
		zap.String("version", app.Version),
		zap.String("commit", app.Commit),
		zap.String("build_at", app.BuildAt),
	)
	return app.rootCmd.ExecuteContext(ctx)
}
