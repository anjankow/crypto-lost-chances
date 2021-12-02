package app

import (
	lostchancescalc "api/internal/lost_chances_calc"
	progressupdates "api/internal/progress_updates"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger         *zap.Logger
	progressReader *progressupdates.Reader
	calcClient     lostchancescalc.Client
}

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string `json:"currency"`
	Income         string `json:"income"`
}

func NewApp(l *zap.Logger, progressReader *progressupdates.Reader) (app App, err error) {
	if progressReader == nil {
		err = errors.New("progress reader is nil")
		return
	}

	app = App{
		Logger:         l,
		progressReader: progressReader,
		calcClient:     lostchancescalc.NewClient(l),
	}
	return
}

func (a App) StartCalculation(ctx context.Context, requestID string, input UserInput) error {
	if err := a.calcClient.StartCalculation(ctx, requestID, input.MonthYear, input.Amount); err != nil {
		return errors.New("calculation request failed: " + err.Error())
	}
	return nil
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

func (a App) GetResults(ctx context.Context, requestID string) (Results, error) {
	return Results{Cryptocurrency: "ADA", Income: "€123"}, nil
}
