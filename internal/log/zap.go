package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOptions struct {
	Named string
}

func NewProductionLogger(opt *LoggerOptions) *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalln("failed init log:", err)
	}

	if opt == nil {
		return logger
	}

	if opt.Named != "" {
		logger = logger.Named("transaction_logger")
	}

	return logger
}
