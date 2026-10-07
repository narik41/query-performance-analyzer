package logger

import (
	"context"
	"sync"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	DebugWithCtx(ctx context.Context, msg string)
	InfoWithCtx(ctx context.Context, msg string)
	WarnWithCtx(ctx context.Context, msg string)
	ErrorWithCtx(ctx context.Context, msg string)
	FatalWithCtx(ctx context.Context, msg string)

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	DebugWithCtxf(ctx context.Context, template string, args ...interface{})
	InfoWithCtxf(ctx context.Context, template string, args ...interface{})
	WarnWithCtxf(ctx context.Context, template string, args ...interface{})
	ErrorWithCtxf(ctx context.Context, template string, args ...interface{})
	FatalWithCtxf(ctx context.Context, template string, args ...interface{})
}

var once sync.Once

var (
	instance Logger
)

func Default() Logger {
	once.Do(func() {
		instance = ZapLogger{}
	})
	return instance
}
