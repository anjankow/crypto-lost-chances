package main

import (
	"api/internal/app"
	"api/internal/config"
	progressupdates "api/internal/progress_updates"
	"api/internal/server"
	"context"
	"log"

	"go.uber.org/zap"
)

func main() {

	logger, err := app.GetLogger()
	if err != nil {
		log.Fatalln("Setting up the logger failed: ", err)
		return
	}
	defer logger.Sync()

	logger.Info("Service started", zap.String("env", string(config.GetRunEnvironment())))

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
	defer service.Close()

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
