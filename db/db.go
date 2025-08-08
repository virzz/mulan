package db

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	zapgorm2 "moul.io/zapgorm2"
)

func NewLogger(isDebug bool) gLogger.Interface {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()
	if isDebug {
		logger.LogMode(gLogger.Info)
	}
	return logger
}

func DefaultConfig() *gorm.Config {
	return &gorm.Config{
		QueryFields:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
		// PrepareStmt:                              true,
	}
}

func New(dialector gorm.Dialector, conn *ConnConfig, opts ...gorm.Option) (*gorm.DB, error) {
	gormCfg := DefaultConfig()
	for _, opt := range opts {
		opt.Apply(gormCfg)
	}

	//Dialector
	// switch DBType(dsnURL.Scheme) {
	// case DBMySQL:
	// 	if strings.HasPrefix(dsnURL.Host, "/") {
	// 		dsnURL.Host = "unix(" + dsnURL.Host + ")"
	// 	} else {
	// 		dsnURL.Host = "tcp(" + dsnURL.Host + ")"
	// 	}
	// 	dsn := dsnURL.String()
	// 	dialector = mysql.New(mysql.Config{
	// 		DSN:                    dsn[strings.Index(dsn, "://")+3:],
	// 		DefaultStringSize:      255,
	// 		DontSupportRenameIndex: true,
	// 	})
	// 	zap.L().Info("Connecting to DB", zap.String("dsn", dsn))
	// case DBPgSQL:
	// 	query := dsnURL.Query()
	// 	if !query.Has("sslmode") {
	// 		query.Set("sslmode", "disable")
	// 	}
	// 	dsnURL.RawQuery = query.Encode()
	// 	dialector = postgres.New(postgres.Config{
	// 		DSN:                  dsnURL.String(),
	// 		PreferSimpleProtocol: true,
	// 	})
	// default:
	// 	return nil, fmt.Errorf("unsupported db type: '%s'", dsnURL.Scheme)
	// }
	// Open
	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		zap.L().Error("Failed to open db", zap.String("name", dialector.Name()), zap.Error(err))
		return nil, err
	}
	// sql.DB Config
	if conn != nil {
		sqlDB, err := db.DB()
		if err != nil {
			zap.L().Warn("Failed to get sql.db", zap.Error(err))
		} else {
			sqlDB.SetMaxIdleConns(conn.Idle)                                     // 最大空闲连接
			sqlDB.SetMaxOpenConns(conn.Open)                                     // 最大连接数
			sqlDB.SetConnMaxLifetime(time.Duration(conn.Lifetime) * time.Second) // 最大可复用时间
		}
	}
	return db, nil
}
