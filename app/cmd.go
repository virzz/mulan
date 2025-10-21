package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func (app *App[T]) internalCmd() {
	app.rootCmd.AddCommand(
		&cobra.Command{
			Use:     "config json|yaml",
			Aliases: []string{"c"},
			Short:   "Show Config Template",
			PreRunE: app.rootCmd.Root().PreRunE,
			RunE: func(cmd *cobra.Command, args []string) error {
				var buf []byte
				if len(args) > 0 && (args[0] == "yaml" || args[0] == "yml") {
					buf, _ = yaml.Marshal(app.conf)
				} else {
					buf, _ = json.MarshalIndent(app.conf, "", "  ")
				}
				fmt.Fprintln(os.Stdout, string(buf))
				return nil
			},
		},

		&cobra.Command{
			Use:     "version",
			Aliases: []string{"v"},
			Short:   "Show Version",
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
				fmt.Fprintln(os.Stdout, buf.String())
			},
		},
	)
}
