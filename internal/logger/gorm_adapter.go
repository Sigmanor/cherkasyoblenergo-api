package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func NewGormLogger(level logger.LogLevel, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		LogLevel:      level,
		SlowThreshold: slowThreshold,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		slog.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		slog.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		slog.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []any{
		slog.Duration("duration", elapsed),
		slog.Int64("rows", rows),
		slog.String("sql", sql),
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fields = append(fields, slog.String("error", err.Error()))
		if l.LogLevel >= logger.Error {
			slog.Error("SQL Error", fields...)
		}
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		if l.LogLevel >= logger.Warn {
			slog.Warn("Slow SQL", fields...)
		}
		return
	}

	if l.LogLevel >= logger.Info {
		slog.Debug("SQL", fields...)
	}
}
