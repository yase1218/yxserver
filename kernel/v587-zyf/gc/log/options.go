package log

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type outputOption struct {
	enable bool
	path   string
	level  zapcore.LevelEnabler
}

type rotateOption struct {
	maxSizeMb      int
	maxBackupFiles int
	maxAgeDays     int
	compress       bool
}

type options struct {
	app string
	env string

	stdOpt   outputOption
	fileOpt  outputOption
	errorOpt outputOption

	rotateOpt rotateOption

	addCaller     bool
	addStacktrace zapcore.LevelEnabler
	development   bool

	sampling *zap.SamplingConfig

	header string
}

type Option func(*options) error

func WithApp(app string) Option {
	return func(c *options) error {
		if app == "" {
			return errors.New("logger app name cannot be empty")
		}
		c.app = app
		return nil
	}
}

func WithStd(enabel bool, level zapcore.LevelEnabler) Option {
	return func(c *options) error {
		c.stdOpt.enable = enabel
		if !c.stdOpt.enable {
			return nil
		}
		c.stdOpt.level = level
		return nil
	}
}

func WithFile(enabel bool, path string, level zapcore.LevelEnabler) Option {
	return func(c *options) error {
		c.fileOpt.enable = enabel
		if !c.fileOpt.enable {
			return nil
		}
		if path != "" {
			c.fileOpt.path = path
		} else {
			c.fileOpt.path = normalPath
		}
		c.fileOpt.level = level
		return nil
	}
}

func WithError(enabel bool, path string, level zapcore.LevelEnabler) Option {
	return func(c *options) error {
		c.errorOpt.enable = enabel
		if !c.errorOpt.enable {
			return nil
		}
		if path != "" {
			c.errorOpt.path = path
		} else {
			c.errorOpt.path = errorPath
		}
		c.errorOpt.level = level
		return nil
	}
}

func WithRotate(maxSizeMb, maxBackupFiles, maxAgeDays int, compress bool) Option {
	return func(c *options) error {
		if maxSizeMb <= 0 {
			return errors.New("logger max size must be positive")
		}
		c.rotateOpt.maxSizeMb = maxSizeMb
		c.rotateOpt.maxBackupFiles = maxBackupFiles
		c.rotateOpt.maxAgeDays = maxAgeDays
		c.rotateOpt.compress = compress
		return nil
	}
}

func WithCaller(enabled bool) Option {
	return func(c *options) error {
		c.addCaller = enabled
		return nil
	}
}

func WithStacktrace(level zapcore.Level) Option {
	return func(c *options) error {
		c.addStacktrace = level
		return nil
	}
}

func WithDevelopment(development bool) Option {
	return func(c *options) error {
		c.development = development
		return nil
	}
}

func WithSampling(initial, thereafter int) Option {
	return func(c *options) error {
		c.sampling = &zap.SamplingConfig{
			Initial:    initial,
			Thereafter: thereafter,
		}
		return nil
	}
}

func WithHeader(header string) Option {
	return func(c *options) error {
		c.header = header
		return nil
	}
}
