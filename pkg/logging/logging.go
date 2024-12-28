package logging

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger
var globalLoggingLevel zapcore.Level

func SetupLogging(_ context.Context) (func() error, error) {
	var err error
	globalLogger, err = zap.NewDevelopment()
	return globalLogger.Sync, err
}

func SetLevel(level zapcore.Level) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	globalLoggingLevel = level
	globalLogger = globalLogger.WithOptions(zap.IncreaseLevel(level))
}

func Info(_ context.Context, msg string, args ...zap.Field) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	globalLogger.Info(msg, args...)
}

func Debug(_ context.Context, msg string, args ...zap.Field) {
	if globalLoggingLevel != zapcore.DebugLevel {
		return
	}
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	globalLogger.Debug(msg, args...)
}

func Error(_ context.Context, msg string, args ...zap.Field) {
	if globalLoggingLevel != zapcore.ErrorLevel {
		return
	}
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	globalLogger.Error(msg, args...)
}
