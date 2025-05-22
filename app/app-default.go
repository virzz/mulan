//go:build !remote
// +build !remote

package app

import (
	"context"

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
		// Config
		instance, _ := cmd.Flags().GetString("instance")
		configPath, _ := cmd.Flags().GetString("config")
		if configPath != "" {
			viper.SetConfigFile(configPath)
		} else {
			viper.SetConfigType("json")
			viper.AddConfigPath(".")
			viper.SetConfigName("config_" + instance)
		}
		err = viper.ReadInConfig()
		if err != nil {
			viper.SetConfigType("yaml")
			err = viper.ReadInConfig()
			if err != nil {
				app.log.Warn("Failed to read in config", zap.Error(err))
			}
		}
		return app.preRunE(cmd.Context())
	}
	return app.rootCmd.ExecuteContext(ctx)
}
