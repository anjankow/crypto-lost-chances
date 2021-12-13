package app

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger() (*zap.Logger, error) {

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.FatalLevel),
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.Development = true
	config.Level.SetLevel(zap.DebugLevel)

	logger, err := config.Build()
	return logger.WithOptions(options...), err

}
