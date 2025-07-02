//go:build !remote
// +build !remote

package app

import (
	"context"
	"log/slog"

	slogzap "github.com/samber/slog-zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Remote struct{}

func (app *App) EnableRemote(project string, publicKey ...string) error {
	panic("remote config build with remote tag")
}

func (app *App) ExecuteE(ctx context.Context) (err error) {
	app.rootCmd.PreRunE = func(cmd *cobra.Command, args []string) (err error) {
		fs := cmd.Flags()
		verbose, _ := fs.GetCount("verbose")
		instance, _ := fs.GetString("instance")
		configPath, _ := fs.GetString("config")
		if verbose > 1 {
			viper.SetOptions(
				viper.WithLogger(
					slog.New(
						slogzap.Option{Level: slog.LevelDebug, Logger: app.log.Named("viper")}.
							NewZapHandler(),
					),
				),
			)
		}
		if configPath != "" {
			viper.SetConfigFile(configPath)
		} else {
			viper.SetConfigType("json")
			viper.AddConfigPath(".")
			viper.SetConfigName("config_" + instance)
		}
		err = viper.ReadInConfig()
		if err != nil {
			app.log.Warn("Failed to read in config", zap.Error(err))
			viper.SetConfigType("yaml")
			err = viper.ReadInConfig()
			if err != nil {
				app.log.Warn("Failed to read in config", zap.Error(err))
			}
		}
		if verbose >= 2 {
			viper.Debug()
		}
		return app.preRunE()
	}

	app.injectVersionCmd()

	if !disableConfigCmd {
		app.injectConfigCmd()
	}
	if app.validate != nil {
		app.injectValidateCmd()
	}
	return app.rootCmd.ExecuteContext(ctx)
}
