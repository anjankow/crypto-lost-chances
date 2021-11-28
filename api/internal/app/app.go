package app

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
}

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string
	Income         float32
}

func NewApp(l *zap.Logger) App {
	return App{
		Logger: l,
	}
}

func (a App) ProcessCalculateReq(ctx context.Context, input UserInput) (Results, error) {
	results := Results{Cryptocurrency: "ADA", Income: float32(input.Amount * 2)}
	return results, nil
}

func (a App) getProgress() {

}
