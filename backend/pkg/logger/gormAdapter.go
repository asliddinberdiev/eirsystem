package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type GormAdapter struct {
	logger        Logger
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

func NewGormAdapter(l Logger, slowThreshold time.Duration) *GormAdapter {
	return &GormAdapter{
		logger:        l,
		SlowThreshold: slowThreshold,
		LogLevel:      gormlogger.Info,
	}
}

func (a *GormAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newAdapter := *a
	newAdapter.LogLevel = level
	return &newAdapter
}

func (a *GormAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if a.LogLevel >= gormlogger.Info {
		a.logger.Info(fmt.Sprintf(msg, data...))
	}
}

func (a *GormAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if a.LogLevel >= gormlogger.Warn {
		a.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (a *GormAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if a.LogLevel >= gormlogger.Error {
		a.logger.Error(fmt.Sprintf(msg, data...))
	}
}

func (a *GormAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if a.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []Field{
		String("sql", sql),
		Duration("duration", elapsed),
		Int64("rows", rows),
	}

	if err != nil && !errors.Is(err, gormlogger.ErrRecordNotFound) {
		fields = append(fields, Error(err))
		a.logger.Error("Database Error", fields...)
		return
	}

	if a.SlowThreshold != 0 && elapsed > a.SlowThreshold {
		a.logger.Warn("Slow SQL Detected", fields...)
		return
	}

	if a.LogLevel == gormlogger.Info {
		a.logger.Info("SQL Query", fields...)
	}
}
