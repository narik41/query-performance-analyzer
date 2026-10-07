package logger

import (
	"context"
	"database-query-analyzer/internal/config"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.SugaredLogger

func init() {
	logger, err := initLogger()
	if err != nil {
		log.Fatalf("Error while initializing zap logger. Error : %s", err)
	} else {
		l = logger
	}
}

type ZapLogger struct {
}

func (z ZapLogger) Debug(args ...interface{}) {
	l.Debug(args...)
}

func (z ZapLogger) Info(args ...interface{}) {
	l.Info(args...)
}

func (z ZapLogger) Warn(args ...interface{}) {
	l.Warn(args...)
}

func (z ZapLogger) Error(args ...interface{}) {
	l.Error(args...)
}

func (z ZapLogger) Fatal(args ...interface{}) {
	l.Fatal(args...)
}

func (z ZapLogger) DebugWithCtx(ctx context.Context, msg string) {
	l.Debug(buildMsgWithCtx(ctx, msg))
}

func (z ZapLogger) InfoWithCtx(ctx context.Context, msg string) {
	l.Info(buildMsgWithCtx(ctx, msg))
}

func (z ZapLogger) WarnWithCtx(ctx context.Context, msg string) {
	l.Warn(buildMsgWithCtx(ctx, msg))
}

func (z ZapLogger) ErrorWithCtx(ctx context.Context, msg string) {
	l.Error(buildMsgWithCtx(ctx, msg))
}

func (z ZapLogger) FatalWithCtx(ctx context.Context, msg string) {
	l.Fatal(buildMsgWithCtx(ctx, msg))
}

func (z ZapLogger) Debugf(template string, args ...interface{}) {
	l.Debugf(template, args...)
}

func (z ZapLogger) Infof(template string, args ...interface{}) {
	l.Infof(template, args...)
}

func (z ZapLogger) Warnf(template string, args ...interface{}) {
	l.Warnf(template, args...)
}

func (z ZapLogger) Errorf(template string, args ...interface{}) {
	l.Errorf(template, args...)
}
func (z ZapLogger) Fatalf(template string, args ...interface{}) {
	l.Fatalf(template, args...)
}

func (z ZapLogger) DebugWithCtxf(ctx context.Context, template string, args ...interface{}) {
	l.Debugf(buildMsgWithCtx(ctx, template), args...)
}

func (z ZapLogger) InfoWithCtxf(ctx context.Context, template string, args ...interface{}) {
	l.Infof(buildMsgWithCtx(ctx, template), args...)
}

func (z ZapLogger) WarnWithCtxf(ctx context.Context, template string, args ...interface{}) {
	l.Warnf(buildMsgWithCtx(ctx, template), args...)
}

func (z ZapLogger) ErrorWithCtxf(ctx context.Context, template string, args ...interface{}) {
	l.Errorf(buildMsgWithCtx(ctx, template), args...)
}

func (z ZapLogger) FatalWithCtxf(ctx context.Context, template string, args ...interface{}) {
	l.Fatalf(buildMsgWithCtx(ctx, template), args...)
}

func getDefaultLogLevel() zapcore.Level {
	level := defaultString(config.Default().GetString("logger.logLevel"), "debug")
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func getConsoleAppender(encoderConfig zapcore.EncoderConfig, defaultLogLevel zapcore.Level) zapcore.Core {
	isConsoleAppender := config.Default().GetBool("logger.consoleAppender")
	if isConsoleAppender {
		consoleAppender := zapcore.NewConsoleEncoder(encoderConfig)
		return zapcore.NewCore(consoleAppender, zapcore.AddSync(os.Stdout), defaultLogLevel)
	}

	return nil
}

func getFileAppender(encoderConfig zapcore.EncoderConfig, defaultLogLevel zapcore.Level) zapcore.Core {
	isFileAppender := config.Default().GetBool("logger.fileAppender")
	if isFileAppender {
		fileAppender := zapcore.NewJSONEncoder(encoderConfig)
		logFileName := defaultString(config.Default().GetString("logger.fileName"), "log.json")
		logFile, _ := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := zapcore.AddSync(logFile)
		return zapcore.NewCore(fileAppender, writer, defaultLogLevel)
	} else {
		return nil
	}
}

func initLogger() (*zap.SugaredLogger, error) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	timeKey := defaultString(config.Default().GetString("logger.timeKey"), "timestamp")
	encoderConfig.TimeKey = timeKey
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.StacktraceKey = "trace"
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	defaultLogLevel := getDefaultLogLevel()
	consoleCore := getConsoleAppender(encoderConfig, defaultLogLevel)
	fileCore := getFileAppender(encoderConfig, defaultLogLevel)
	var cores []zapcore.Core
	if consoleCore != nil {
		cores = append(cores, consoleCore)
	}
	if fileCore != nil {
		cores = append(cores, fileCore)
	}

	if len(cores) > 0 {
		core := zapcore.NewTee(cores...)
		logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
		return logger.Sugar(), nil
	}

	return nil, errors.New("could not create any logger")
}

func defaultString(str string, defaultStr string) string {
	if str == "" {
		return defaultStr
	}

	return str
}

func buildMsgWithCtx(ctx context.Context, msg string) string {
	traceKey := GetTraceKey()
	return fmt.Sprintf("%s %s : %s", msg, traceKey, ctx.Value(traceKey))
}
