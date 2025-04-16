package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"

	"github.com/virzz/mulan/utils/once"
	"github.com/virzz/vlog"
)

var (
	std      *gorm.DB
	oncePlus once.OncePlus
	multi    cmap.ConcurrentMap[string, *gorm.DB]
)

func R() *gorm.DB {
	if std == nil {
		panic("db not init")
	}
	return std
}

func Migrate(models ...any) error { return std.AutoMigrate(models...) }

func connect(cfg *Config) (*gorm.DB, error) {
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
	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		vlog.Error("Failed to connect db", "err", err.Error())
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		vlog.Warn("Failed to get sql.db", "err", err.Error())
	} else {
		sqlDB.SetMaxIdleConns(cfg.Conn.Idle)                                     // 最大空闲连接
		sqlDB.SetMaxOpenConns(cfg.Conn.Open)                                     // 最大连接数
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Conn.Lifetime) * time.Second) // 最大可复用时间
	}
	if cfg.Debug {
		db = db.Debug()
	}
	return db, nil
}

func Init(cfg *Config, force ...bool) error {
	if len(force) > 0 && force[0] {
		db, err := connect(cfg)
		if err != nil {
			return err
		}
		std = db
		return nil
	}
	return oncePlus.Do(func() (err error) {
		db, err := connect(cfg)
		if err != nil {
			return err
		}
		std = db
		return nil
	})
}

func New(cfg *Config, name string) (*gorm.DB, error) {
	if name == "" {
		name = cfg.Name
	}
	if multi.Has(name) {
		return nil, fmt.Errorf("db %s already exists", name)
	}
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	multi.Set(name, db)
	return db, nil
}
