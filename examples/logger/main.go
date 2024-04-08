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
	logOpt = log.CommonLogOpt.WithDirectory("/tmp/1234").WithTraceIDEnable(false).WithLogLevel(zapcore.DebugLevel)
)

func main() {
	logger := log.GetLogger("example", &logOpt)

	logger.Debug(ctx, "debug log")
	logger.Info(ctx, "info log")
	logger.Warn(ctx, "warn log")
	logger.Error(ctx, "error log")

	logger.Info(traceIDCtx, "info log with traceid")
}
