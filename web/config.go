package web

import (
	"strconv"

	"github.com/spf13/pflag"
)

func FlagSet(defaultPort int) *pflag.FlagSet {
	fs := pflag.NewFlagSet("http", pflag.ContinueOnError)
	fs.String("http.system", "", "HTTP System Token")
	fs.String("http.prefix", "/api", "HTTP API Route Prefix")
	fs.String("http.endpoint", "", "HTTP Domain Endpoint")
	fs.String("http.host", "127.0.0.1", "HTTP Listen Address")
	fs.Int("http.port", defaultPort, "HTTP Listen Port")
	fs.Bool("http.debug", false, "HTTP Debug Mode")
	fs.Bool("http.pprof", false, "Enable PProf")
	fs.Bool("http.requestid", false, "Enable HTTP RequestID")
	fs.Bool("http.metrics", false, "Enable Metrics")

	fs.Bool("token.disable", false, "Token Disable")
	fs.String("token.system", "", "Token System Token")
	fs.String("token.keyname", "", "Token APIKey Name")
	fs.String("token.apikey", "", "Token APIKey")
	return fs
}

type (
	Token struct {
		Disable bool   `json:"disable" yaml:"disable"`
		KeyName string `json:"keyname" yaml:"keyname"`
		APIKey  string `json:"apikey" yaml:"apikey"`
		System  string `json:"system" yaml:"system"`
	}

	//go:generate structx -struct Config
	Config struct {
		Prefix    string `json:"prefix" yaml:"prefix" default:"/api"`
		Endpoint  string `json:"endpoint" yaml:"endpoint"`
		Host      string `json:"host" yaml:"host" default:"127.0.0.1"`
		Port      int    `json:"port" yaml:"port" default:"8080"`
		Debug     bool   `json:"debug" yaml:"debug"`
		Pprof     bool   `json:"pprof" yaml:"pprof"`
		RequestID bool   `json:"requestid" yaml:"requestid"`
		Metrics   bool   `json:"metrics" yaml:"metrics"`
	}
)

func (c *Config) Addr() string { return c.Host + ":" + strconv.Itoa(c.Port) }

func (c *Config) GetEndpoint() string {
	if c.Endpoint == "" {
		c.Endpoint = c.Addr()
	}
	return c.Endpoint
}

func (c *Config) Check() *Config {
	if c.Endpoint == "" {
		c.Endpoint = c.Addr()
	}
	if c.Prefix == "" {
		c.Prefix = "/api"
	}
	if c.Port == 0 {
		c.Port = 5678
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	return c
}
