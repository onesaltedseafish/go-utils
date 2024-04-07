// Package gorm Logger Implementation
package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/onesaltedseafish/go-utils/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

// Logger logger for gorm
type Logger struct {
	level zapcore.Level
	log   *log.Logger
}

// NewLogger new logger
func NewLogger(name string, opt *log.LoggerOpt) *Logger {
	l := Logger{}
	l.log = log.GetLogger(name, opt)
	return &l
}

// LogMode log mode
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	switch level {
	case logger.Info:
		newLogger.level = zapcore.InfoLevel
	case logger.Warn:
		newLogger.level = zapcore.WarnLevel
	case logger.Error:
		newLogger.level = zapcore.ErrorLevel
	default:
		newLogger.level = zapcore.DebugLevel
	}
	newLogger.log.SetLevel(newLogger.level)
	return &newLogger
}

// Info print info
func (l *Logger) Info(ctx context.Context, msg string, _ ...interface{}) {
	l.log.Info(ctx, msg)
}

// Warn print warn messages
func (l *Logger) Warn(ctx context.Context, msg string, _ ...interface{}) {
	l.log.Warn(ctx, msg)
}

// Error print error messages
func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warn(ctx, msg)
}

// Trace print sql message
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	l.log.Debug(
		ctx,
		"trace sql message",
		zap.String("elapsed_time", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds()/1e6))),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Error(err),
	)
}
