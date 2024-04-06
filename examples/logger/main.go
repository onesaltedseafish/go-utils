// Package main TODO
package main

import (
	"context"

	"github.com/onesaltedseafish/go-utils/log"
	"go.uber.org/zap/zapcore"
)

var (
	ctx        = context.Background()
	traceIDCtx = log.NewTraceIdWithCtx(ctx)
)

var (
	logOpt = log.LoggerOpt{
		LogLevel:      zapcore.DebugLevel,
		Directory:     "/tmp/1234",
		TraceIDEnable: true,
		MaxSize:       10,
		MaxBackups:    10,
		MaxAge:        1,
	}
)

func main() {
	logger := log.GetLogger("example", &logOpt)

	logger.Debug(ctx, "debug log")
	logger.Info(ctx, "info log")
	logger.Warn(ctx, "warn log")
	logger.Error(ctx, "error log")

	logger.Info(traceIDCtx, "info log with traceid")
}
