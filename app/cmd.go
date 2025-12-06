package app

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func (app *App[T]) internalCmd() {
	app.rootCmd.AddCommand(
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
