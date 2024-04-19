// Package main TODO
package main

import (
	"context"

	"github.com/onesaltedseafish/go-utils/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ctx        = context.Background()
	traceIDCtx = log.NewTraceIdWithCtx(ctx)
)

var (
	logOpt          = log.CommonLogOpt.WithDirectory("/tmp/1234").WithTraceIDEnable(false).WithLogLevel(zapcore.DebugLevel)
	noConsoleLogOpt = log.CommonLogOpt.WithDirectory("/tmp/1234").WithTraceIDEnable(false).
			WithLogLevel(zapcore.DebugLevel).
			WithConsoleLog(false)
)

func main() {
	logger := log.GetLogger("example", &logOpt)

	logger2 := log.NewFromLogger(logger).WithLoggerMetaFields(zap.String("version", "2"))

	logger3 := log.GetLogger("example3", &noConsoleLogOpt)

	logger.Debug(ctx, "debug log")
	logger.Info(ctx, "info log")
	logger2.Info(ctx, "info log 2")
	logger3.Info(ctx, "test for no console log")
	logger.Warn(ctx, "warn log")
	logger.Error(ctx, "error log")

	logger.Info(traceIDCtx, "info log with traceid")
	logger.Fatal(ctx, "fatal log")
}
