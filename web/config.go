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
	fs.Bool("http.auth", false, "Enable Auth")
	fs.Bool("http.debug", false, "HTTP Debug Mode")
	fs.Bool("http.pprof", false, "Enable PProf")
	fs.Bool("http.requestid", false, "Enable HTTP RequestID")
	fs.Bool("http.metrics", false, "Enable Metrics")
	fs.StringSlice("http.origins", []string{"*"}, "HTTP CORS: Allow Origins")
	fs.StringSlice("http.headers", []string{"Authorization"}, "HTTP CORS: Allow Headers")
	return fs
}

//go:generate structx -struct Config
type Config struct {
	System    string   `json:"system" yaml:"system"`
	Prefix    string   `json:"prefix" yaml:"prefix" default:"/api"`
	Endpoint  string   `json:"endpoint" yaml:"endpoint"`
	Host      string   `json:"host" yaml:"host" default:"127.0.0.1"`
	Port      int      `json:"port" yaml:"port" default:"8080"`
	Origins   []string `json:"origins" yaml:"origins"`
	Headers   []string `json:"headers" yaml:"headers"`
	Debug     bool     `json:"debug" yaml:"debug"`
	Pprof     bool     `json:"pprof" yaml:"pprof"`
	RequestID bool     `json:"requestid" yaml:"requestid"`
	Metrics   bool     `json:"metrics" yaml:"metrics"`
	Auth      bool     `json:"auth" yaml:"auth"`
}

func (c *Config) GetEndpoint() string {
	if c.Endpoint == "" {
		c.Endpoint = c.Addr()
	}
	return c.Endpoint
}
func (c *Config) Addr() string { return c.Host + ":" + strconv.Itoa(c.Port) }
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
