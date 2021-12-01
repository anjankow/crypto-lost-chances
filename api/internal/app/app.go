package app

import (
	progressupdates "api/internal/progress_updates"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger         *zap.Logger
	progressReader *progressupdates.Reader
}

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string
	Income         float32
}

func NewApp(l *zap.Logger, progressReader *progressupdates.Reader) (app App, err error) {
	if progressReader == nil {
		err = errors.New("progress reader is nil")
		return
	}

	app = App{
		Logger:         l,
		progressReader: progressReader,
	}
	return
}

func (a App) StartCalculation(ctx context.Context, input UserInput) (Results, error) {
	results := Results{Cryptocurrency: "ADA", Income: float32(input.Amount * 2)}

	// calls the main app

	return results, nil
}

// ListenProgress listens on the queue for the request progress
func (a App) ListenProgress(ctx context.Context, requestID string, callback func(progress int)) {
	channel, finish := a.progressReader.SubscribeToProgressUpdates(requestID)
	defer finish()

	for p := range channel {
		a.Logger.Debug("received a progress update", zap.String("requestID", requestID))
		callback(p)
	}

	// progress == 100, here possibly some other actions on this event

}
