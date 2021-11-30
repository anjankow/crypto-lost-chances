package main

import (
	"api/internal/app"
	progressupdates "api/internal/progress_updates"
	"api/internal/server"
	"context"
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

	reader := progressupdates.NewReader(logger)
	readerCloser, err := reader.Start(context.Background())
	defer readerCloser()

	if err != nil {
		logger.Fatal("failed to start the progress updates reader: " + err.Error())
	}

	service, err := app.NewApp(logger, &reader)
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
