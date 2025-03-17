package db

import (
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/virzz/vlog"
)

func DialectorPgSQL(cfg *Config) gorm.Dialector {
	dsnList := []string{}
	if cfg.Host != "" {
		dsnList = append(dsnList, "host="+cfg.Host)
	}
	if cfg.Port != 0 {
		dsnList = append(dsnList, "port="+strconv.Itoa(cfg.Port))
	}
	if cfg.User != "" {
		dsnList = append(dsnList, "user="+cfg.User)
	}
	if cfg.Pass != "" {
		dsnList = append(dsnList, "password="+cfg.Pass)
	}
	if cfg.Name != "" {
		dsnList = append(dsnList, "dbname="+cfg.Name)
	}
	if len(dsnList) > 0 {
		dsnList = append(dsnList, "sslmode=disable", "TimeZone=Asia/Shanghai")
	}
	dsn := strings.Join(dsnList, " ")
	if cfg.Debug {
		vlog.Info("Connecting to postgres", "dsn", dsn)
	}
	return postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true})
}
