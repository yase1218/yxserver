package log

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *Logger
	once     sync.Once
)

func Instance() *Logger {
	return instance
}

func Init(opts ...Option) {
	var err error
	once.Do(func() {
		instance, err = New(opts...)
	})
	if err != nil {
		panic(fmt.Errorf("logger init failed,err: %w", err))
	}
}

func Debug(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Error(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Panic(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Fatal(msg, fields...)
	}
}

func Debugf(format string, args ...interface{}) {
	if instance != nil {
		instance.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if instance != nil {
		instance.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if instance != nil {
		instance.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if instance != nil {
		instance.Errorf(format, args...)
	}
}

func Panicf(format string, args ...interface{}) {
	if instance != nil {
		instance.Panicf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if instance != nil {
		instance.Fatalf(format, args...)
	}
}

func With(fields ...zap.Field) *zap.Logger {
	if instance != nil {
		return instance.With(fields...)
	}
	return nil
}

func Sync() error {
	if instance != nil {
		return instance.Sync()
	}
	return nil
}

// ==============================
// Convenient global initialization function
// ==============================

func InitDevelopment(app, head string) {
	Init(WithApp(app),
		WithStd(true, zapcore.DebugLevel),
		WithFile(true, "", zapcore.DebugLevel),
		WithError(true, "", zapcore.ErrorLevel),
		WithCaller(true),
		WithStacktrace(zapcore.DPanicLevel),
		WithHeader(head),
	)
}

// InitProduction
func InitProduction(app, head string) {
	Init(WithApp(app),
		WithStd(false, zapcore.DebugLevel),
		WithFile(true, "", zapcore.InfoLevel),
		WithError(true, "", zapcore.ErrorLevel),
		WithCaller(true),
		WithSampling(10, 100),
		WithHeader(head),
	)
}
