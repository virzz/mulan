package app

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/virzz/mulan/db"
	"github.com/virzz/mulan/rdb"
	"github.com/virzz/mulan/web"
)

type (
	Configer interface {
		Validate() error
		GetHTTP() *web.Config
		GetDB() *db.Config
	}
	Config struct {
		HTTP web.Config `json:"http" yaml:"http"`
		DB   db.Config  `json:"db" yaml:"db"`
		RDB  rdb.Config `json:"rdb" yaml:"rdb"`
	}
)

func (c *Config) GetHTTP() *web.Config { return &c.HTTP }
func (c *Config) GetDB() *db.Config    { return &c.DB }
func (c *Config) GetRDB() *rdb.Config  { return &c.RDB }
func (c *Config) Validate() error      { return nil }
func (c *Config) Template(typ ...string) string {
	_c := &Config{
		HTTP: web.Config{
			Prefix:    "/api",
			Port:      5678,
			Host:      "0.0.0.0",
			RequestID: true,
			Metrics:   true,
			Auth:      true,
		},
		DB: db.Config{
			Host: "127.0.0.1",
			Type: "postgres",
			Port: 5432,
			User: "postgres",
			Pass: "postgres",
			Name: "postgres",
		},
		RDB: rdb.Config{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}
	_type := "json"
	if len(typ) > 0 {
		_type = typ[0]
	}
	var (
		buf []byte
		err error
	)
	switch _type {
	case "yaml":
		buf, err = yaml.Marshal(&_c)
	case "json":
		fallthrough
	default:
		buf, err = json.MarshalIndent(&_c, "", "  ")
	}
	if err != nil {
		return err.Error()
	}
	return string(buf)
}
