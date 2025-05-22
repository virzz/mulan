package app

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/virzz/mulan/db"
	"gopkg.in/yaml.v3"
)

func (app *App) AddMaintain() {
	app.rootCmd.AddGroup(&cobra.Group{ID: "maintain", Title: "Maintain Commands"})
	app.rootCmd.AddCommand(db.MaintainCommand(app.conf.GetDB())...)
}

func (app *App) AddConfigCmd() {
	if !slices.ContainsFunc(app.rootCmd.Commands(),
		func(cmd *cobra.Command) bool { return cmd.Use == "config" },
	) {
		app.rootCmd.AddCommand(&cobra.Command{
			Use: "config json|yaml", Aliases: []string{"c"},
			Short: "Show Config Template",
			RunE: func(cmd *cobra.Command, args []string) error {
				var buf []byte
				var config any
				if app.conf != nil {
					config = app.conf
				} else {
					config = viper.AllSettings()
					viper.Set("config", nil)
					viper.Set("instance", nil)
				}
				if len(args) > 0 && (args[0] == "yaml" || args[0] == "yml") {
					buf, _ = yaml.Marshal(config)
				} else {
					buf, _ = json.MarshalIndent(config, "", "  ")
				}
				fmt.Println(string(buf))
				return nil
			},
		})
	}
}

func (app *App) AddValidateCmd() {
	app.rootCmd.AddCommand(&cobra.Command{
		Use: "validate", Aliases: []string{"valid"},
		Short: "Validate Configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.conf.Validate(); err != nil {
				return err
			}
			if app.validate != nil {
				if err := app.validate(); err != nil {
					return err
				}
			}
			fmt.Println("Valid configuration")
			return nil
		},
	})
}
