// Package logger - Zap Logger
package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/asliddinberdiev/eirsystem/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Field = zapcore.Field

var (
	Int      = zap.Int
	Int64    = zap.Int64
	String   = zap.String
	Error    = zap.Error
	Bool     = zap.Bool
	Any      = zap.Any
	Duration = zap.Duration
	Time     = zap.Time
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	Named(name string) Logger
	With(fields ...Field) Logger

	Sync() error
}

type loggerImpl struct {
	zap *zap.Logger
}

func New(cfg config.Logger, isDev bool) (Logger, error) {
	if cfg.LogDir != "" {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			return nil, err
		}
	}

	atomicLevel := zap.NewAtomicLevel()
	parsedLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		parsedLevel = zapcore.InfoLevel
	}
	atomicLevel.SetLevel(parsedLevel)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	var cores []zapcore.Core

	if cfg.LogDir != "" && cfg.Filename != "" {
		logPath := filepath.Join(cfg.LogDir, cfg.Filename)

		rotator := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}

		if cfg.RotateDaily {
			go func() {
				for {
					now := time.Now()
					nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

					time.Sleep(time.Until(nextMidnight))

					_ = rotator.Rotate()
				}
			}()
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(rotator),
			atomicLevel,
		)
		cores = append(cores, fileCore)
	}
	if cfg.Console {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.Lock(os.Stdout),
			atomicLevel,
		)
		cores = append(cores, consoleCore)
	}

	core := zapcore.NewTee(cores...)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	if isDev {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	zapLogger := zap.New(core, options...)

	return &loggerImpl{
		zap: zapLogger,
	}, nil
}

func (l *loggerImpl) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

func (l *loggerImpl) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

func (l *loggerImpl) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *loggerImpl) Named(name string) Logger {
	return &loggerImpl{
		zap: l.zap.Named(name),
	}
}

func (l *loggerImpl) With(fields ...Field) Logger {
	return &loggerImpl{
		zap: l.zap.With(fields...),
	}
}

func (l *loggerImpl) Sync() error {
	return l.zap.Sync()
}
