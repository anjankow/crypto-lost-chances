package app

import "go.uber.org/zap"

type App struct {
	Logger *zap.Logger
}

func NewApp(l *zap.Logger) App {
	return App{
		Logger: l,
	}
}
