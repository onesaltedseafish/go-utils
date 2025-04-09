package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceID(t *testing.T) {
	testcases := []struct {
		ctx  context.Context
		Want string
	}{
		{context.TODO(), ""},
		{context.WithValue(context.TODO(), ctxTraceIdKey, "1234"), "1234"},
	}

	for _, testcase := range testcases {
		assert.Equal(t, testcase.Want, GetTraceIDWithCtx(testcase.ctx))
	}
}

func TestNewTraceIDCtx(t *testing.T) {
	ctx := NewTraceIDWithCtx(context.Background())
	if GetTraceIDWithCtx(ctx) == "" {
		t.Errorf("NewTraceIdWithCtx didn't generate traceid")
	}
}
