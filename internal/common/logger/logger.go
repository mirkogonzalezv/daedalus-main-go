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

		if env == "" {
			env = "development"
		}

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
			fallback := zap.NewExample()
			fallback.Error("No se pudo inicializar zap correctamente, usando fallback", zap.Error(err))
			log = fallback
			return
		}

		log = l
	})
}

func L() *zap.Logger {
	if log == nil {
		panic("logger no inicializado — llama a logger.Init(env) antes de usar logger.L()")
	}
	return log
}

func Sugar() *zap.SugaredLogger {
	return L().Sugar()
}
