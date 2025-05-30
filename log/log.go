package log

import (
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Kafka struct {
		Broker   []string `json:"broker" yaml:"broker"`
		Topic    string   `json:"topic" yaml:"topic"`
		Username string   `json:"username" yaml:"username"`
		Password string   `json:"password" yaml:"password"`
		Level    string   `json:"level" yaml:"level"`
	}
	Http struct {
		URL   string `json:"url" yaml:"url"`
		Level string `json:"level" yaml:"level"`
	}
	Config struct {
		IsDev bool    `json:"is_dev" yaml:"is_dev"`
		Level string  `json:"level" yaml:"level"`
		Kafka []Kafka `json:"kafka,omitzero" yaml:"kafka"`
		Http  []Http  `json:"http,omitzero" yaml:"http"`
		File  string  `json:"file,omitzero" yaml:"file"`
	}
)

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("log", pflag.ContinueOnError)
	fs.Bool("log.dev", false, "use DevelopmentEncoder")
	fs.String("log.level", "info", "log level")
	fs.String("log.file", "", "log file")
	return fs
}

func New(lvl zapcore.Level, isDev bool, name ...string) (*zap.Logger, error) {
	cfg := &Config{Level: lvl.String(), IsDev: isDev}
	return NewWithConfig(cfg, name...)
}

func NewWithConfig(cfg *Config, name ...string) (*zap.Logger, error) {
	var encoder zapcore.Encoder
	if cfg.IsDev {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		_prodCfg := zap.NewProductionEncoderConfig()
		_prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(_prodCfg)
	}
	lvl, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lvl),
	}
	var prodCfg = zap.NewProductionEncoderConfig()
	for _, h := range cfg.Http {
		lvl, err := zapcore.ParseLevel(h.Level)
		if err != nil {
			lvl = zapcore.InfoLevel
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(prodCfg),
			newHTTPWriter(h.URL),
			lvl,
		))
	}
	if cfg.File != "" {
		err = os.MkdirAll(filepath.Dir(cfg.File), 0755)
		if err != nil {
			return nil, err
		}
		f, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(prodCfg),
			zapcore.Lock(f),
			lvl,
		))
	}
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.DPanicLevel),
	)
	if len(name) > 0 {
		logger = logger.Named(name[0]).
			With(zap.String("service", name[0]))
	}
	zap.ReplaceGlobals(logger)
	return logger, nil
}
