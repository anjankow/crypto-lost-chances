package main

import (
	"context"
	"log"
	"lost-chances-calc/internal/app"
	progressupdates "lost-chances-calc/internal/progress_updates"
	"lost-chances-calc/internal/server"
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

	writer := progressupdates.NewWriter(logger)
	close, err := writer.Init(context.Background())
	defer close()

	if err != nil {
		logger.Fatal("failed to initialize the pubsub writer: " + err.Error())
	}

	service, err := app.NewApp(logger, &writer)
	if err != nil {
		logger.Fatal("service creation failed: " + err.Error())
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
