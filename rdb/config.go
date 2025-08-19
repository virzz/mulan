package rdb

import (
	"strconv"

	"github.com/spf13/pflag"
)

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("rdb", pflag.ContinueOnError)
	fs.Bool("rdb.debug", false, "Database Debug Mode")
	fs.String("rdb.host", "127.0.0.1", "Database Host")
	fs.Int("rdb.port", 6379, "Database Port")
	fs.Int("rdb.db", 0, "Database Index")
	fs.String("rdb.pass", "", "Database Password")
	return fs
}

type Config struct {
	Debug bool   `json:"debug" yaml:"debug"`
	Host  string `json:"host" yaml:"host"`
	Port  int    `json:"port" yaml:"port"`
	DB    int    `json:"db" yaml:"db"`
	Pass  string `json:"pass" yaml:"pass"`
}

func (c *Config) FlagSet() *pflag.FlagSet { return FlagSet() }

func (c *Config) Addr() string { return c.Host + ":" + strconv.Itoa(c.Port) }

func (c *Config) GetRDB() *Config { return c }
