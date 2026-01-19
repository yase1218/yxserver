package db_model

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DbLogger struct {
}

func (dl *DbLogger) Printf(format string, v ...any) {
	anys := make([]zapcore.Field, len(v))
	for _, a := range v {
		anys = append(anys, zap.Any("-", a))
	}

	log.Info(format, anys...)
}
