package main

import (
	"api/internal/app"
	"api/internal/server"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	logger, err := getLogger()
	if err != nil {
		log.Fatalln("Setting up the logger failed: ", err)
		return
	}
	defer logger.Sync()

	logger.Info("Service started")

	service := app.NewApp(logger)
	if err != nil {
		logger.Fatal("New service creation failed: ", zap.Error(err))
		return
	}

	// HTTP SERVER
	ser := server.NewServer(logger, &service)
	if err != nil {
		logger.Fatal("Server creation failed: ", zap.Error(err))
	}

	err = ser.Run()

	if err != nil {
		logger.Error("Service finished with error", zap.Error(err))
	} else {
		logger.Info("Service finished")
	}
}

func getLogger() (*zap.Logger, error) {
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
