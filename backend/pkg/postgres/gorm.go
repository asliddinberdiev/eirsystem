package postgres

import (
	"context"
	"time"

	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

func New(cfg *config.Postgres, isDev bool, log logger.Logger) (*gorm.DB, error) {
	adapter := logger.NewGormAdapter(log, 200*time.Millisecond)

	var mode gormLog.LogLevel = gormLog.Error
	if isDev {
		mode = gormLog.Info
	}

	newLogger := adapter.LogMode(mode)

	db, err := gorm.Open(
		postgres.Open(cfg.GetDSN()),
		&gorm.Config{
			Logger:                 newLogger,
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			QueryFields:            true,
		},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}