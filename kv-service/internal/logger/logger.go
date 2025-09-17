package logger

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.CallerKey = ""
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalln(err)
	}
	return logger.Sugar()
}
