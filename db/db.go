package db

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	zapgorm2 "moul.io/zapgorm2"

	"github.com/virzz/mulan/utils/once"
)

var (
	std      *gorm.DB
	oncePlus once.OncePlus
	multi    = cmap.New[*gorm.DB]()
)

func R(name ...string) *gorm.DB {
	if len(name) > 0 {
		if db, ok := multi.Get(name[0]); ok {
			return db
		}
		panic(name[0] + " not init")
	}
	if std == nil {
		panic("db not init")
	}
	return std
}

func Migrate(models ...any) error { return std.AutoMigrate(models...) }

func connect(cfg *Config, wrapper ...*DialectorWrapper) (*gorm.DB, error) {
	dsnURL, err := url.Parse(cfg.DSN)
	if err != nil {
		zap.L().Error("parse dsn fail:", zap.Error(err))
		return nil, err
	}
	if cfg.User != "" || cfg.Pass != "" {
		dsnURL.User = url.UserPassword(cfg.User, cfg.Pass)
	}
	if dsnURL.Host == "" {
		dsnURL.Host = "localhost"
	}
	var newLogger gLogger.Interface
	if cfg.Debug {
		newLogger = gLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	} else {
		logger := zapgorm2.New(zap.L())
		logger.SetAsDefault()
		newLogger = logger
	}
	gormCfg := &gorm.Config{
		Logger:                                   newLogger,
		QueryFields:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	}
	if cfg.DisablePrepareStmt {
		gormCfg.PrepareStmt = false
	}
	//Dialector
	var dialector gorm.Dialector
	switch DBType(dsnURL.Scheme) {
	case DBMySQL:
		if strings.HasPrefix(dsnURL.Host, "/") {
			dsnURL.Host = "unix(" + dsnURL.Host + ")"
		} else {
			dsnURL.Host = "tcp(" + dsnURL.Host + ")"
		}
		query := dsnURL.Query()
		if !query.Has("charset") {
			query.Set("charset", "utf8mb4")
		}
		if !query.Has("parseTime") {
			query.Set("parseTime", "True")
		}
		if !query.Has("loc") {
			query.Set("loc", "Local")
		}
		dsnURL.RawQuery = query.Encode()
		dsn := dsnURL.String()
		dialector = mysql.New(mysql.Config{
			DSN:                    dsn[strings.Index(dsn, "://")+3:],
			DefaultStringSize:      255,
			DontSupportRenameIndex: true,
		})
		zap.L().Info("Connecting to DB", zap.String("dsn", dsn))
	case DBPgSQL:
		fallthrough
	default:
		query := dsnURL.Query()
		if !query.Has("sslmode") {
			query.Set("sslmode", "disable")
		}
		if !query.Has("TimeZone") {
			query.Set("TimeZone", "Asia/Shanghai")
		}
		dsnURL.RawQuery = query.Encode()
		dialector = postgres.New(postgres.Config{
			DSN:                  dsnURL.String(),
			PreferSimpleProtocol: true,
		})
	}
	// Open
	if len(wrapper) > 0 {
		wrapper[0].Apply(dialector)
	}
	db, err := gorm.Open(wrapper[0], gormCfg)
	if err != nil {
		zap.L().Error("Failed to connect db", zap.String("dsn", dsnURL.String()), zap.Error(err))
		return nil, err
	}
	// sql.DB Config
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Warn("Failed to get sql.db", zap.Error(err))
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

func New(cfg *Config, name string, wrapper ...*DialectorWrapper) (*gorm.DB, error) {
	if name == "" {
		name = "std"
	}
	if multi.Has(name) {
		return nil, fmt.Errorf("db %s already exists", name)
	}
	db, err := connect(cfg, wrapper...)
	if err != nil {
		return nil, err
	}
	multi.Set(name, db)
	return db, nil
}
