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
		DSN:                       dsn,
		DefaultStringSize:         255,  // string 类型字段的默认长度
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
	})
}
