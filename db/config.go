package db

import (
	"github.com/spf13/pflag"
)

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("db", pflag.ContinueOnError)
	fs.Bool("db.debug", false, "Database Debug Mode")
	fs.String("db.dsn", "", "Database DSN")
	fs.String("db.user", "", "Database User")
	fs.String("db.pass", "", "Database Password")
	fs.String("db.name", "", "Database Name")
	fs.Int("db.conn.idle", 20, "Database MaxIdleConns")
	fs.Int("db.conn.open", 250, "Database MaxOpenConns")
	fs.Int("db.conn.lifetime", 3600, "Database ConnMaxLifetime")
	return fs
}

type (
	ConnConfig struct {
		Idle     int `json:"idle,omitempty" yaml:"idle,omitempty"`
		Open     int `json:"open,omitempty" yaml:"open,omitempty"`
		Lifetime int `json:"lifetime,omitempty" yaml:"lifetime,omitempty"`
	}
	Config struct {
		Debug              bool        `json:"debug,omitempty" yaml:"debug,omitempty"`
		DSN                string      `json:"dsn,omitempty" yaml:"dsn,omitempty"`
		User               string      `json:"user,omitempty" yaml:"user,omitempty"`
		Pass               string      `json:"pass,omitempty" yaml:"pass,omitempty"`
		Name               string      `json:"name,omitempty" yaml:"name,omitempty"`
		Conn               *ConnConfig `json:"conn,omitempty" yaml:"conn,omitempty"`
		DisablePrepareStmt bool        `json:"disable_prepare_stmt,omitempty" yaml:"disable_prepare_stmt,omitempty"`
	}
	DBType string
)

const (
	DBMySQL DBType = "mysql"
	DBPgSQL DBType = "postgres"
)
