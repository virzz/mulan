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
