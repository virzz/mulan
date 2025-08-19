package db

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func FlagSet(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.Bool(name+".debug", false, "Database Debug Mode")
	fs.String(name+".dsn", "", "Database DSN")
	fs.String(name+".user", "", "Database User")
	fs.String(name+".pass", "", "Database Password")
	fs.String(name+".name", "", "Database Name")
	fs.Int(name+".conn.idle", 20, "Database MaxIdleConns")
	fs.Int(name+".conn.open", 250, "Database MaxOpenConns")
	fs.Int(name+".conn.lifetime", 3600, "Database ConnMaxLifetime")
	return fs
}

type (
	ConnConfig struct {
		Idle     int `json:"idle,omitempty" yaml:"idle,omitempty"`
		Open     int `json:"open,omitempty" yaml:"open,omitempty"`
		Lifetime int `json:"lifetime,omitempty" yaml:"lifetime,omitempty"`
	}
	Config struct {
		Debug bool              `json:"debug,omitempty" yaml:"debug,omitempty"`
		DSN   string            `json:"dsn,omitempty" yaml:"dsn,omitempty"`
		User  string            `json:"user,omitempty" yaml:"user,omitempty"`
		Pass  string            `json:"pass,omitempty" yaml:"pass,omitempty"`
		Name  string            `json:"name,omitempty" yaml:"name,omitempty"`
		Conn  *ConnConfig       `json:"conn,omitempty" yaml:"conn,omitempty"`
		Args  map[string]string `json:"args,omitempty" yaml:"args,omitempty"`
	}
)

func (c *Config) FlagSet(name string) *pflag.FlagSet { return FlagSet(name) }

func (c *Config) WithArgs(k, v string) *Config {
	if c.Args == nil {
		c.Args = make(map[string]string)
	}
	c.Args[k] = v
	return c
}

func (c *Config) String() string {
	dsnURL, err := url.Parse(c.DSN)
	if err != nil {
		zap.L().Error("parse dsn fail:", zap.Error(err))
		return ""
	}
	_user := dsnURL.User.Username()
	_pass, _ := dsnURL.User.Password()
	if c.User != "" {
		_user = c.User
	}
	if c.Pass != "" {
		_pass = c.Pass
	}
	if _user != "" || _pass != "" {
		dsnURL.User = url.UserPassword(_user, _pass)
	}
	if c.Name != "" {
		dsnURL.Path = "/" + c.Name
	}
	if c.Args != nil {
		query := dsnURL.Query()
		for k, v := range c.Args {
			query.Set(k, v)
		}
		dsnURL.RawQuery = query.Encode()
	}

	dsn := dsnURL.String()
	switch dsnURL.Scheme {
	case "mysql":
		if strings.HasPrefix(dsnURL.Host, "/") {
			dsnURL.Host = "unix(" + dsnURL.Host + ")"
		} else {
			if !strings.Contains(dsnURL.Host, "tcp") {
				dsnURL.Host = "tcp(" + dsnURL.Host + ")"
			}
		}
		dsn = dsnURL.String()
		dsn, _ = strings.CutPrefix(dsn, "mysql://")
	case "postgres":
	case "sqlite3":
		fmt.Println("host", dsnURL.Host, "path", dsnURL.Path)
		if dsnURL.Host == ":memory:" {
			dsn = "file::memory:"
		} else {
			dsn = filepath.Join(".", dsnURL.Host, dsnURL.Path)
		}
	}
	return dsn
}

func (c *Config) GetDB() *Config { return c }
