package app

import (
	"context"
	"errors"
	progressupdates "lost-chances-calc/internal/progress_updates"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger         *zap.Logger
	progressWriter *progressupdates.Writer
}

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string `json:"currency"`
	Income         string `json:"income"`
}

func NewApp(l *zap.Logger, progressWriter *progressupdates.Writer) (app App, err error) {
	if progressWriter == nil {
		err = errors.New("progress reader is nil")
		return
	}

	app = App{
		Logger:         l,
		progressWriter: progressWriter,
	}
	return
}

func (a App) StartCalculation(ctx context.Context, input UserInput) error {

	// calls the main app

	return nil
}
