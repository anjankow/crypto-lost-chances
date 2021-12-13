package app

import (
	lostchancescalc "api/internal/lost_chances_calc"
	progressupdates "api/internal/progress_updates"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	maxProgress = 100
)

type App struct {
	Logger         *zap.Logger
	progressReader *progressupdates.Reader
	calcClient     lostchancescalc.Client

	results *map[string](chan calcRequestResult)
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string `json:"currency"`
	Income         string `json:"income"`
}

type calcRequestResult struct {
	CalcResults Results
	Error       error
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
		results:        &map[string]chan calcRequestResult{},
		mutex:          &sync.Mutex{},
		wg:             &sync.WaitGroup{},
	}
	return
}

func (a *App) Close() {
	a.wg.Wait()
}

func (a App) StartCalculation(ctx context.Context, requestID string, input UserInput) error {

	(*a.results)[requestID] = make(chan calcRequestResult)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		lostChance, err := a.calcClient.Calculate(ctx, requestID, input.MonthYear, input.Amount)
		if err != nil {
			(*a.results)[requestID] <- calcRequestResult{Error: err}

			return
		}

		res := convertLostChanceToApiResults(lostChance)

		(*a.results)[requestID] <- calcRequestResult{Error: nil, CalcResults: res}
		close((*a.results)[requestID])

	}()

	return nil
}

func (a App) GetResults(ctx context.Context, requestID string) (Results, error) {
	r := <-(*a.results)[requestID]
	a.Logger.Debug("received calc result from the channel", zap.String("requestID", requestID))

	if r.Error != nil {
		return Results{}, r.Error
	}

	return r.CalcResults, nil
}

// ListenProgress listens on the queue for the request progress
func (a App) ListenProgress(ctx context.Context, requestID string, callback func(progress int)) {
	channel := a.progressReader.Subscribe(requestID)

	for p := range channel {
		a.Logger.Debug("received a progress update", zap.String("requestID", requestID))
		callback(p)

		if p >= maxProgress {
			a.Logger.Debug("progress reached max, unsubscribing", zap.String("requestID", requestID))
			a.progressReader.Unsubscribe(requestID)
			break
		}
	}

}

func convertLostChanceToApiResults(chance lostchancescalc.LostChance) (r Results) {
	r.Cryptocurrency = chance.CryptocurName
	r.Income = fmt.Sprintf("€ %v", chance.OutputFiatAmount-chance.InputFiatAmount)
	return r
}
