package db

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/virzz/vlog"
)

func DialectorMySQL(cfg *Config) gorm.Dialector {
	addr := cfg.Host
	if addr == "" {
		addr = "localhost"
	}
	if strings.HasPrefix(addr, "/") {
		addr = "unix(" + addr + ")"
	} else if cfg.Port != 0 {
		addr = "tcp(" + addr + ":" + strconv.Itoa(cfg.Port) + ")"
	} else {
		addr = "tcp(" + addr + ")"
	}
	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	if cfg.Name == "" {
		cfg.Name = "mysql"
	}
	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, addr, cfg.Name, charset,
	)
	if cfg.Debug {
		vlog.Info("Connecting to MySQL", "dsn", dsn)
	}
	return mysql.New(mysql.Config{
		DSN:                    dsn,
		DefaultStringSize:      255,
		DontSupportRenameIndex: true,
	})
}
