package log

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	app = "app"
	env = "dev"

	normalPath = "./log/info/"
	errorPath  = "./log/error/"

	normalSuffix = ".log"
	errorSuffix  = ".err.log"

	defaultMaxSizeMb      = 100
	defaultMaxBackupFiles = 100
	defaultMaxAgeDays     = 90

	defaultSamplingInitial    = 100
	defaultSamplingThereafter = 100
)

type Logger struct {
	opts  *options
	base  *zap.Logger
	sugar *zap.SugaredLogger

	header string
}

func New(vals ...Option) (*Logger, error) {
	opts := defaultOptions()
	for _, opt := range vals {
		if err := opt(opts); err != nil {
			return nil, fmt.Errorf("logger invalid option: %w", err)
		}
	}

	if err := buildPath(opts); err != nil {
		return nil, fmt.Errorf("logger create path failed : %w", err)
	}

	core, err := buildCore(opts)
	if err != nil {
		return nil, fmt.Errorf("logger build core failed: %w", err)
	}

	zapOptions := buildZapOptions(opts)
	zapLogger := zap.New(core, zapOptions...)
	return &Logger{
		base:   zapLogger,
		sugar:  zapLogger.Sugar(),
		opts:   opts,
		header: "-" + opts.header + "- ",
	}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field)  { l.base.Info(l.header+msg, fields...) }
func (l *Logger) Debug(msg string, fields ...zap.Field) { l.base.Debug(l.header+msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.base.Error(l.header+msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.base.Warn(l.header+msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.base.Fatal(l.header+msg, fields...) }
func (l *Logger) Panic(msg string, fields ...zap.Field) { l.base.Panic(l.header+msg, fields...) }

func (l *Logger) Infof(format string, args ...interface{})  { l.sugar.Infof(l.header+format, args...) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.sugar.Debugf(l.header+format, args...) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.sugar.Errorf(l.header+format, args...) }
func (l *Logger) Warnf(format string, args ...interface{})  { l.sugar.Warnf(l.header+format, args...) }
func (l *Logger) Fatalf(format string, args ...interface{}) { l.sugar.Fatalf(l.header+format, args...) }
func (l *Logger) Panicf(format string, args ...interface{}) { l.sugar.Panicf(l.header+format, args...) }

func (l *Logger) Write(p []byte) (n int, err error) {
	l.base.Info(string(p))
	return len(p), nil
}

func (l *Logger) With(fields ...zap.Field) *zap.Logger { return l.base.With(fields...) }

func (l *Logger) Sync() error { return l.base.Sync() }
