package app

import (
	"context"
	"errors"
	"lost-chances-calc/internal/domain"
	pricefetcher "lost-chances-calc/internal/price_fetcher"
	progressupdates "lost-chances-calc/internal/progress_updates"
	"sync"
	"time"

	"go.uber.org/zap"
)

const priceFetcherResultsTimeout = 5 * time.Second
const progressMax = 100

var (
	progressSteps = 2*len(domain.Currencies) + //for submitting fetcher tasks and fetching the results
		1 + // for the initial progress added just because the operation reached this point
		2 // for getting the current price and calculating the final result

	progressStepLen = int(progressMax / progressSteps)
)

type App struct {
	Logger             *zap.Logger
	progressWriter     *progressupdates.Writer
	priceFetchListener *pricefetcher.Listener
	priceFetcher       pricefetcher.PriceFetcher

	wg *sync.WaitGroup
}

type CalcInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string `json:"currency"`
	Income         string `json:"income"`
}

func NewApp(l *zap.Logger, progressWriter *progressupdates.Writer) (app App, err error) {
	if progressWriter == nil {
		err = errors.New("progress writer is nil")
		return
	}

	app = App{
		Logger:         l,
		progressWriter: progressWriter,
		wg:             &sync.WaitGroup{},
	}
	return
}

func (a App) Calculate(ctx context.Context, requestID string, input CalcInput) error {

	a.priceFetchListener.Register(ctx, requestID)
	defer a.priceFetchListener.Unregister(requestID)

	a.progressWriter.PublishProgress(ctx, requestID, progress)

	historicalPrices, err := a.getHistoricalPrices(ctx, requestID, &progress, domain.Euro, input.MonthYear) // currently only for euro
	if err != nil {
		a.Logger.Error("getting the historical price failed: "+err.Error(), zap.String("requestID", requestID))
		return err
	}

	currentPrices, err := a.getCurrentPrices(ctx, domain.Euro)
	if err != nil {
		a.Logger.Error("getting the current prices failed: "+err.Error(), zap.String("requestID", requestID))
		return err
	}

	return

	// a.wg.Add(1)
	// go func() {
	// 	defer a.wg.Done()

	// 	ctx := context.Background() // in not in the request context

	// 	for i := 0; i <= 100; i += 20 {
	// 		a.progressWriter.PublishProgress(ctx, requestID, i)
	// 		time.Sleep(1 * time.Second)
	// 	}

	// }()

	return nil
}

func (a App) publishProgress(ctx context.Context, currentProgress *int, requestID string) {

	*currentProgress += progressStepLen
	if *currentProgress > progressMax {
		*currentProgress = progressMax
	}
	if progressMax-*currentProgress < progressStepLen {
		*currentProgress = progressMax
	}

	a.progressWriter.PublishProgress(ctx, requestID, *currentProgress)
}
