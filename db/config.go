package db

import (
	"os/user"

	"github.com/spf13/pflag"
)

func FlagSet(name string) *pflag.FlagSet {
	username := "root"
	u, _ := user.Current()
	if u != nil {
		username = u.Username
	}
	fs := pflag.NewFlagSet("db", pflag.ContinueOnError)
	fs.Bool("db.debug", false, "Database Debug Mode")
	fs.String("db.host", "127.0.0.1", "Database Host/UnixSocket")
	fs.Int("db.port", 5432, "Database Port")
	fs.String("db.name", name, "Database Name")
	fs.String("db.user", username, "Database User")
	fs.String("db.pass", "", "Database Password")
	fs.Int("db.conn.idle", 20, "Database MaxIdleConns")
	fs.Int("db.conn.open", 250, "Database MaxOpenConns")
	fs.Int("db.conn.lifetime", 3600, "Database ConnMaxLifetime")
	return fs
}

type ConnConfig struct {
	Idle     int `json:"idle" yaml:"idle"`
	Open     int `json:"open" yaml:"open"`
	Lifetime int `json:"lifetime" yaml:"lifetime"`
}

type DBType string

const (
	DBMySQL DBType = "mysql"
	DBPgSQL DBType = "postgres"
)

const (
	defaultQueryPath string = "./dao/query"
)

type Config struct {
	Debug              bool       `json:"debug" yaml:"debug"`
	DSN                string     `json:"dsn" yaml:"dsn"`
	User               string     `json:"user" yaml:"user"`
	Pass               string     `json:"pass" yaml:"pass"`
	Conn               ConnConfig `json:"conn" yaml:"conn"`
	DisablePrepareStmt bool       `json:"disable_prepare_stmt" yaml:"disable_prepare_stmt"`
}
