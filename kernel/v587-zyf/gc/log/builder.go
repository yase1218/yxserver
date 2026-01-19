package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func defaultOptions() *options {
	return &options{
		app: app,
		env: env,
		stdOpt: outputOption{
			enable: true,
			level:  zapcore.InfoLevel,
		},
		fileOpt: outputOption{
			enable: true,
			level:  zapcore.InfoLevel,
		},
		errorOpt: outputOption{
			enable: true,
			level:  zapcore.ErrorLevel,
		},

		rotateOpt: rotateOption{
			maxSizeMb:      defaultMaxSizeMb,
			maxBackupFiles: defaultMaxBackupFiles,
			maxAgeDays:     defaultMaxAgeDays,
			compress:       true,
		},

		sampling: &zap.SamplingConfig{
			Initial:    defaultSamplingInitial,
			Thereafter: defaultSamplingThereafter,
		},
	}
}

func buildPath(cfg *options) error {
	paths := make(map[string]bool)

	if cfg.fileOpt.enable {
		paths[cfg.fileOpt.path] = true
	}

	if cfg.errorOpt.enable {
		paths[cfg.errorOpt.path] = true
	}

	for path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("create log dir %s failed: %w", path, err)
		}
	}
	return nil
}

func buildCore(cfg *options) (zapcore.Core, error) {
	var cores []zapcore.Core
	encoderCfg := zapcore.EncoderConfig{
		// some of the following fileds are used only by json encoder
		TimeKey:       "tm",
		LevelKey:      "lvl",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",

		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
	}
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	if cfg.fileOpt.enable {
		var fileLevel zapcore.LevelEnabler
		if cfg.development {
			fileLevel = zapcore.DebugLevel
		} else {
			fileLevel = cfg.fileOpt.level
		}
		cores = append(cores,
			zapcore.NewCore(encoder, buildWriter(cfg.fileOpt.path+cfg.app+normalSuffix, &cfg.rotateOpt), fileLevel))
	}

	if cfg.errorOpt.enable {
		cores = append(cores,
			zapcore.NewCore(encoder, buildWriter(cfg.errorOpt.path+cfg.app+errorSuffix, &cfg.rotateOpt), cfg.errorOpt.level))
	}

	if cfg.stdOpt.enable {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		stdEncoder := zapcore.NewConsoleEncoder(encoderCfg)
		var stdLevel zapcore.LevelEnabler
		if cfg.development {
			stdLevel = zapcore.DebugLevel
		} else {
			stdLevel = cfg.stdOpt.level
		}
		cores = append(cores, zapcore.NewCore(stdEncoder, zapcore.AddSync(colorable.NewColorableStdout()), stdLevel))
	}

	if len(cores) == 0 {
		return nil, fmt.Errorf("logger no log output enabled")
	}
	return zapcore.NewTee(cores...), nil
}

func buildZapOptions(cfg *options) []zap.Option {
	var opts []zap.Option
	if cfg.addCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(2))
	}
	if cfg.addStacktrace != nil {
		opts = append(opts, zap.AddStacktrace(cfg.addStacktrace))
	}
	if cfg.development {
		opts = append(opts, zap.Development())
	}
	if cfg.sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				cfg.sampling.Initial,
				cfg.sampling.Thereafter,
			)
		}))
	}
	return opts
}

func buildWriter(fileName string, cfg *rotateOption) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    cfg.maxSizeMb,
		MaxBackups: cfg.maxBackupFiles,
		MaxAge:     cfg.maxAgeDays,
		LocalTime:  true,
		Compress:   cfg.compress,
	})
}
