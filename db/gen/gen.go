package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"

	"github.com/virzz/mulan/db"
	gMysql "github.com/virzz/mulan/db/gen/mysql"
	mlog "github.com/virzz/mulan/log"
)

var rootCmd = &cobra.Command{
	Use:          "",
	Short:        "mulan-gen is a tool for generating gorm models",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config, _ := cmd.Flags().GetString("config")
		if config != "" {
			viper.SetConfigFile(config)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}
		return nil
	},
	RunE: genProcess,
}

func genProcess(cmd *cobra.Command, args []string) error {
	var cfg = &db.GenConfig{}
	if err := viper.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return err
	}
	config := cfg.Database.Revise()
	wrapper := db.NewDialectorWrapper()
	wrapper.SetMigrator(func(tx *gorm.DB) gorm.Migrator {
		dialector := tx.Dialector.(*db.DialectorWrapper).Dialector.(*mysql.Dialector)
		return gMysql.Migrator{
			Migrator: migrator.Migrator{
				Config: migrator.Config{
					DB:        tx,
					Dialector: dialector,
				},
			},
			Dialector: *dialector,
		}
	})
	db, err := db.New(&db.Config{DSN: config.DSN, DisablePrepareStmt: true}, "std", wrapper)
	if err != nil {
		return err
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:           config.OutPath,
		OutFile:           config.OutFile,
		ModelPkgPath:      config.ModelPkgName,
		WithUnitTest:      config.WithUnitTest,
		FieldNullable:     config.FieldNullable,
		FieldCoverable:    config.FieldCoverable,
		FieldWithIndexTag: config.FieldWithIndexTag,
		FieldWithTypeTag:  config.FieldWithTypeTag,
		FieldSignable:     config.FieldSignable,
	})
	g.UseDB(db)
	// Generate models
	if len(config.Tables) == 0 {
		config.Tables, err = db.Migrator().GetTables()
		if err != nil {
			return fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}
	models := make([]any, len(config.Tables))
	for i, tableName := range config.Tables {
		models[i] = g.GenerateModel(tableName)
	}
	if !config.OnlyModel {
		g.ApplyBasic(models...)
	}
	g.Execute()
	return nil
}

func main() {
	mlog.NewWithConfig(mlog.Config{Level: "error"})
	rootCmd.Flags().AddFlagSet(db.GenFlagSet())
	viper.BindPFlags(rootCmd.Flags())
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("parse config fail:", zap.Error(err))
	}
}
