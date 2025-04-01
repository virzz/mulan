package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"

	"github.com/virzz/mulan/utils/once"
	"github.com/virzz/vlog"
)

var (
	std      *gorm.DB
	oncePlus once.OncePlus
)

func R() *gorm.DB {
	if std == nil {
		panic("db not init")
	}
	return std
}

func Migrate(models ...any) error { return std.AutoMigrate(models...) }

func Init(cfg *Config, force ...bool) error {
	if len(force) > 0 && force[0] {
		return connect(cfg)
	}
	return oncePlus.Do(func() (err error) {
		return connect(cfg)
	})
}

func connect(cfg *Config) (err error) {
	newLogger := gLogger.Default.LogMode(gLogger.Info)
	if cfg.Debug {
		newLogger = gLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	} else {
		f, err := os.OpenFile(filepath.Join("logs", "db.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			vlog.Warn("Failed to open gorm log file", "err", err.Error())
		} else {
			newLogger = gLogger.New(log.New(f, "\r\n", log.LstdFlags),
				gLogger.Config{LogLevel: gLogger.Warn, IgnoreRecordNotFoundError: true},
			)
		}
	}
	gormCfg := &gorm.Config{Logger: newLogger,
		QueryFields:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	}
	//Dialector
	var dialector gorm.Dialector
	switch cfg.Type {
	case DBMySQL:
		dialector = DialectorMySQL(cfg)
	case DBPgSQL:
		fallthrough
	default:
		dialector = DialectorPgSQL(cfg)
	}
	// Open
	std, err = gorm.Open(dialector, gormCfg)
	if err != nil {
		vlog.Error("Failed to connect db", "err", err.Error())
		return err
	}
	sqlDB, err := std.DB()
	if err != nil {
		vlog.Warn("Failed to get sql.db", "err", err.Error())
	} else {
		sqlDB.SetMaxIdleConns(cfg.Conn.Idle)                                     // 最大空闲连接
		sqlDB.SetMaxOpenConns(cfg.Conn.Open)                                     // 最大连接数
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Conn.Lifetime) * time.Second) // 最大可复用时间
	}
	if cfg.Debug {
		std = std.Debug()
	}
	return nil
}
