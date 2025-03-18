package db

import (
	"bytes"
	"fmt"
	"iter"
	"os/user"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	dropPgSQL = `DROP DATABASE IF EXISTS {{.Name}};
DO $$ BEGIN IF EXISTS ( SELECT 1 FROM pg_roles WHERE rolname = '{{.User}}' ) THEN REVOKE ALL PRIVILEGES ON SCHEMA public FROM {{.User}};REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM {{.User}};REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM {{.User}};REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public FROM {{.User}}; END IF ; END $$ ;
DROP USER IF EXISTS {{.User}};
DROP ROLE IF EXISTS {{.User}};`
	createPgSQL = `CREATE ROLE {{.User}} WITH LOGIN PASSWORD '{{.Pass}}';
CREATE DATABASE {{.Name}} OWNER {{.User}};
GRANT ALL PRIVILEGES ON DATABASE {{.Name}} TO {{.User}};
GRANT ALL PRIVILEGES ON SCHEMA public TO {{.User}};
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO {{.User}};
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO {{.User}};
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO {{.User}};`

	dropMySQL   = "DROP DATABASE IF EXISTS `{{.Name}}`;\nDROP USER '{{.User}}'@'%';"
	createMySQL = "CREATE USER IF NOT EXISTS '{{.User}}'@'%' identified by '{{.Pass}}';\nCREATE DATABASE IF NOT EXISTS `{{.Name}}` CHARACTER SET = 'utf8mb4' COLLATE = 'utf8mb4_unicode_ci';\nGRANT ALL ON `{{.Name}}`.* to '{{.User}}'@'%';\nFLUSH PRIVILEGES;"

	hostPgSQL = `127.0.0.1:5432;/tmp;/var/run/postgresql/`
	hostMySQL = `127.0.0.1:3306;/tmp/mysql.sock;/run/mysqld/mysqld.sock;/var/run/mysqld/mysqld.sock;/var/lib/mysql/mysql.sock`
)

func tryConnect(cmd *cobra.Command) (err error) {
	cfg := &Config{Debug: true}
	cfg.Name, _ = cmd.Flags().GetString("name")
	cfg.User, _ = cmd.Flags().GetString("user")
	cfg.Pass, _ = cmd.Flags().GetString("pass")
	cfg.Host, _ = cmd.Flags().GetString("host")
	cfg.Port, _ = cmd.Flags().GetInt("port")
	_type, _ := cmd.Flags().GetString("type")
	cfg.Type = DBType(_type)
	// 自定义 Host 连接
	if cfg.Host != "" {
		if err = Init(cfg, true); err == nil {
			return nil
		}
	}
	var hosts iter.Seq[string]
	switch cfg.Type {
	case DBMySQL:
		hosts = strings.SplitSeq(hostMySQL, ";")
	case DBPgSQL:
		fallthrough
	default:
		hosts = strings.SplitSeq(hostPgSQL, ";")
	}
	cfg.Port = 0
	for host := range hosts {
		cfg.Host = host
		if err = Init(cfg, true); err == nil {
			return nil
		}
	}
	return errors.New("failed to connect database")
}

func parseSql(content string, data any) (string, error) {
	buf := bytes.Buffer{}
	t, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func MaintainCommand(cfg *Config) []*cobra.Command {
	createCmd := &cobra.Command{
		GroupID: "maintain",
		Use:     "create",
		Short:   "Create database (Create User,Database and Grant Privileges)",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tryConnect(cmd)
			if err != nil {
				return err
			}
			_type, _ := cmd.Flags().GetString("type")
			cfg.Type = DBType(_type)
			force, _ := cmd.Flags().GetBool("force")
			var dropSql, createSql string
			switch cfg.Type {
			case DBMySQL:
				dropSql = dropMySQL
				createSql = createMySQL
			case DBPgSQL:
				fallthrough
			default:
				dropSql = dropPgSQL
				createSql = createPgSQL
			}
			if force {
				result, err := parseSql(dropSql, *cfg)
				if err != nil {
					return err
				}
				for sql := range strings.SplitSeq(result, "\n") {
					if err = R().Exec(sql).Error; err != nil {
						return err
					}
				}
			}
			result, err := parseSql(createSql, *cfg)
			if err != nil {
				return err
			}
			for sql := range strings.SplitSeq(result, "\n") {
				fmt.Println(sql)
				if err = R().Exec(sql).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
	username := "root"
	u, _ := user.Current()
	if u != nil {
		username = u.Username
	}
	createCmd.Flags().StringP("type", "T", "", "Database Type: mysql/postgres")
	createCmd.Flags().StringP("host", "H", "", "Database Host/Socket")
	createCmd.Flags().IntP("port", "p", 0, "Database Port")
	createCmd.Flags().StringP("user", "U", username, "Database username")
	createCmd.Flags().StringP("pass", "P", "", "Database  password")
	createCmd.Flags().StringP("name", "N", "", "Database name")
	createCmd.Flags().BoolP("force", "F", false, "Force to delete if exists")

	migrateCmd := &cobra.Command{
		GroupID: "maintain",
		Use:     "migrate",
		Short:   "Migrate database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Init(cfg); err != nil {
				return err
			}
			return Migrate(Models()...)
		},
	}

	initCmd := &cobra.Command{
		GroupID: "maintain",
		Use:     "init",
		Short:   "Init database default data",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = Init(cfg); err != nil {
				return err
			}
			force, _ := cmd.Flags().GetBool("force")
			for _, m := range Models() {
				if m, ok := m.(Modeler); ok {
					if force {
						if err = R().Exec("TRUNCATE TABLE " + m.TableName() + ";").Error; err != nil {
							return err
						}
					}
					// Insert Or Update
					for _, item := range m.Defaults() {
						key, val := item.Unique()
						if err = R().Model(item).
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
	initCmd.Flags().BoolP("force", "F", false, "Force to delete if exists")
	return []*cobra.Command{
		createCmd,
		migrateCmd,
		initCmd,
	}
}
