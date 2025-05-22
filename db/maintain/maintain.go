package maintain

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/virzz/mulan/db"
	"go.uber.org/zap"
)

func flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.BoolP("debug", "D", false, "Debug mode")
	fs.StringP("type", "T", "", "Database Type")
	fs.StringP("schema", "S", "public", "Database Schema (default: public)")
	fs.StringP("host", "H", "", "Database Host/Socket")
	fs.StringP("user", "U", "", "Database username")
	fs.StringP("pass", "P", "", "Database  password")
	fs.StringP("name", "N", "", "Database name")
	fs.String("dsn", "", "Database DSN (URL)")
	return fs
}

func getConfigDSNURL(cmd *cobra.Command) (dsnURL *url.URL, err error) {
	dsnURL = &url.URL{}
	dsn, _ := cmd.Flags().GetString("dsn")
	if dsn != "" {
		dsnURL, err = url.Parse(dsn)
		if err != nil {
			return nil, err
		}
	} else {
		dsnURL.Scheme, _ = cmd.Flags().GetString("type")
		dsnURL.Path, _ = cmd.Flags().GetString("name")
		dsnURL.Host, _ = cmd.Flags().GetString("host")
	}
	if dsnURL.Host == "" {
		return nil, errors.New("database host(host:port/socket) or dsn is required")
	}
	if dsnURL.Scheme == "" {
		return nil, errors.New("database type or dsn is required")
	}
	if dsnURL.Path == "" {
		return nil, errors.New("database name or dsn is required")
	}

	// Auth
	user, _ := cmd.Flags().GetString("user")
	pass, _ := cmd.Flags().GetString("pass")
	_user := dsnURL.User.Username()
	_pass, _ := dsnURL.User.Password()
	if user != "" {
		_user = user
	}
	if pass != "" {
		_pass = pass
	}
	dsnURL.User = url.UserPassword(_user, _pass)
	// Query
	query := dsnURL.Query()
	if dsnURL.Scheme == string(db.DBPgSQL) {
		schema, _ := cmd.Flags().GetString("schema")
		if schema != "" {
			query.Set("search_path", schema)
		}
	}
	dsnURL.RawQuery = query.Encode()
	return dsnURL, nil
}

func Command(dbCfg *db.Config) []*cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create database (Create User,Database and Grant Privileges)",
		RunE: func(cmd *cobra.Command, args []string) error {
			preRunE := cmd.Root().PreRunE
			if preRunE != nil {
				if err := preRunE(cmd, args); err != nil {
					return err
				}
			}
			if dbCfg.DSN == "" {
				return errors.New("database dsn is required")
			}
			dsnURL, err := getConfigDSNURL(cmd)
			if err != nil {
				return err
			}
			debug, _ := cmd.Flags().GetBool("debug")
			dsn := dsnURL.String()
			cfg := &db.Config{Debug: debug, DSN: dsn}
			if err = db.Init(cfg, true); err != nil {
				zap.L().Error("failed to connect", zap.String("dsn", dsn), zap.Error(err))
				return err
			}
			var dropSql, createSql string
			switch db.DBType(dsnURL.Scheme) {
			case db.DBMySQL:
				dropSql = dropMySQL
				createSql = createMySQL
			case db.DBPgSQL:
				dropSql = dropPgSQL
				createSql = createPgSQL
			default:
				return errors.New("unsupported database type")
			}
			data := &SqlTemplateData{
				User:   dbCfg.User,
				Pass:   dbCfg.Pass,
				Name:   strings.TrimPrefix(dsnURL.Path, "/"),
				Schema: dsnURL.Query().Get("search_path"),
			}
			force, _ := cmd.Flags().GetBool("force")
			if force {
				result, err := parseSql(dropSql, *data)
				if err != nil {
					return err
				}
				for sql := range strings.SplitSeq(result, "\n") {
					if err = db.R().Exec(sql).Error; err != nil {
						return err
					}
				}
			}
			result, err := parseSql(createSql, *data)
			if err != nil {
				zap.L().Error("parse create sql fail", zap.Error(err))
				return err
			}
			for sql := range strings.SplitSeq(result, "\n") {
				zap.L().Info("create", zap.String("sql", sql))
				if err = db.R().Exec(sql).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		RunE: func(cmd *cobra.Command, args []string) error {
			preRunE := cmd.Root().PreRunE
			if preRunE != nil {
				if err := preRunE(cmd, args); err != nil {
					return err
				}
			}
			fmt.Printf("migrateCmd %+v\n", dbCfg)
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				dbCfg.Debug = true
			}
			if err := db.Init(dbCfg); err != nil {
				return err
			}
			return db.Migrate(Models()...)
		},
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init database default data",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = db.Init(dbCfg); err != nil {
				return err
			}
			force, _ := cmd.Flags().GetBool("force")
			for _, m := range Models() {
				if m, ok := m.(db.Modeler); ok {
					if force {
						err = db.R().Exec("TRUNCATE TABLE " + m.TableName() + ";").Error
						if err != nil {
							return err
						}
					}
					// Insert Or Update
					for _, item := range m.Defaults() {
						key, val := item.Unique()
						if err = db.R().Model(item).
							Where(key+" = ?", val).
							Attrs(&item).
							FirstOrCreate(item).Error; err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
	}

	createCmd.Flags().AddFlagSet(flagSet())
	createCmd.Flags().BoolP("force", "F", false, "Force to delete if exists")

	initCmd.Flags().BoolP("force", "F", false, "Force to delete if exists")

	return []*cobra.Command{
		createCmd,
		migrateCmd,
		initCmd,
	}
}
