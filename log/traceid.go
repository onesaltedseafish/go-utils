package log

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var (
	ctxTraceIdKey = logCtxTraceIdKey{}
)

type logCtxTraceIdKey struct{}

// NewTraceIDWithCtx wrap ctx with traceid
func NewTraceIDWithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxTraceIdKey, uuid.New().String())
}

// GetTraceIDWithCtx get trace id from ctx
func GetTraceIDWithCtx(ctx context.Context) string {
	value := ctx.Value(ctxTraceIdKey)
	if value != nil {
		return fmt.Sprintf("%s", value)
	}
	// get traceid from otel span
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		return spanContext.TraceID().String()
	}
	return ""
}
