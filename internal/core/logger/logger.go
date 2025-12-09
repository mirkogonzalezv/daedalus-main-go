package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

func Init(env string) {
	once.Do(func() {
		var cfg zap.Config
		if env == "prod" || env == "production" {
			cfg = zap.NewProductionConfig()
			cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		}

		l, err := cfg.Build()
		if err != nil {
			panic("No se puede inicializar zap logger: " + err.Error())
		}

		log = l
	})
}

func L() *zap.Logger {
	if log == nil {
		panic("logger no puede inicializar el llamado logger.Init() en main.go")
	}
	return log
}

// Sugar expone una versión simplificada del logger
func Sugar() *zap.SugaredLogger {
	return L().Sugar()
}
