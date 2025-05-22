package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/virzz/mulan/db"
	"github.com/virzz/mulan/db/maintain"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	disableConfigCmd = false
)

func (app *App) injectMaintainCmd(name string, dbCfg *db.Config) {
	cmd := &cobra.Command{GroupID: "maintain", Use: "m-" + name}
	cmd.AddCommand(maintain.Command(dbCfg)...)
	app.rootCmd.AddCommand(cmd)
}

func (app *App) injectConfigCmd() {
	if slices.ContainsFunc(app.rootCmd.Commands(),
		func(cmd *cobra.Command) bool { return cmd.Use == "config" },
	) {
		zap.L().Panic("config command already exists")
		return
	}
	app.rootCmd.AddCommand(&cobra.Command{
		Use: "config json|yaml", Aliases: []string{"c"},
		Short: "Show Config Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			preRunE := cmd.Root().PreRunE
			if preRunE != nil {
				if err := preRunE(cmd, args); err != nil {
					return err
				}
			}
			var config any
			if app.conf != nil {
				config = app.conf
			} else {
				config = viper.AllSettings()
				viper.Set("config", nil)
				viper.Set("instance", nil)
			}
			var buf []byte
			if len(args) > 0 && (args[0] == "yaml" || args[0] == "yml") {
				buf, _ = yaml.Marshal(config)
			} else {
				buf, _ = json.MarshalIndent(config, "", "  ")
			}
			fmt.Fprintln(os.Stdout, string(buf))
			return nil
		},
	})
}

func (app *App) injectValidateCmd() {
	if slices.ContainsFunc(app.rootCmd.Commands(),
		func(cmd *cobra.Command) bool { return cmd.Use == "validate" },
	) {
		zap.L().Panic("validate command already exists")
		return
	}
	app.rootCmd.AddCommand(&cobra.Command{
		Use: "validate", Aliases: []string{"valid"},
		Short: "Validate Configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			preRunE := cmd.Root().PreRunE
			if preRunE != nil {
				if err := preRunE(cmd, args); err != nil {
					return err
				}
			}
			if err := app.conf.Validate(); err != nil {
				return err
			}
			if app.validate != nil {
				if err := app.validate(); err != nil {
					return err
				}
			}
			zap.L().Info("Valid configuration")
			return nil
		},
	})
}

func (app *App) injectVersionCmd() {
	app.rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show Version",
		Run: func(cmd *cobra.Command, args []string) {
			if app.Commit == "" || app.BuildAt == "" {
				bi, ok := debug.ReadBuildInfo()
				if ok {
					for _, setting := range bi.Settings {
						if setting.Key == "vcs.revision" && app.Commit == "" {
							app.Commit = setting.Value
						}
						if setting.Key == "vcs.time" && app.BuildAt == "" {
							app.BuildAt = setting.Value
						}
					}
				}
			}
			buf := bytes.NewBuffer(nil)
			buf.WriteString(app.Name)
			if app.ID != "" {
				buf.WriteString("(")
				buf.WriteString(app.ID)
				buf.WriteString(")")
			}
			buf.WriteString(" - ")
			buf.WriteString(app.Version)
			if app.Commit != "" {
				buf.WriteString("-")
				buf.WriteString(app.Commit)
			}
			if app.BuildAt != "" {
				buf.WriteString("-")
				buf.WriteString(app.BuildAt)
			}
			if app.Description != "" {
				buf.WriteString(" - ")
				buf.WriteString(app.Description)
			}
			fmt.Println(buf.String())
		},
	})
}
