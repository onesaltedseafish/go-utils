package log

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

var (
	ctxTraceIdKey = logCtxTraceIdKey{}
)

type logCtxTraceIdKey struct{}

// NewTraceIdWithCtx TODO
func NewTraceIdWithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxTraceIdKey, uuid.New().String())
}

// GetTraceIdWithCtx get trace id from ctx
func GetTraceIdWithCtx(ctx context.Context) string {
	value := ctx.Value(ctxTraceIdKey)
	if value == nil {
		return ""
	} else {
		return fmt.Sprintf("%s", value)
	}
}
