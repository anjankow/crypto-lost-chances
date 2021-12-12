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

const (
	priceFetcherResultsTimeout = 5 * time.Second
	progressMax                = 100
	// support only euro for now
	fiatName = domain.Euro
)

var (
	progressSteps = 2*len(domain.Cryptocurrencies) + //for submitting fetcher tasks and fetching the results
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
	Amount    float64
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

func (a App) Calculate(ctx context.Context, requestID string, input CalcInput) (chance *domain.LostChance, err error) {

	progress := progressStepLen
	// if we are here, then some work has been done already - update the progress
	a.progressWriter.PublishProgress(ctx, requestID, progress)

	// request to dispatch the getter tasks and subscribe for the historical price messages
	getHistoricalPrices, err := a.requestHistoricalPrices(ctx, requestID, &progress, fiatName, input.MonthYear)
	if err != nil {
		a.Logger.Error("requesting the historical prices failed: "+err.Error(), zap.String("requestID", requestID))
	}

	currentPrices, err := a.getCurrentPrices(ctx, requestID, &progress, fiatName)
	if err != nil {
		a.Logger.Error("getting the current prices failed: "+err.Error(), zap.String("requestID", requestID))
		return nil, err
	}

	cctx, cancel := context.WithTimeout(ctx, priceFetcherResultsTimeout)
	defer cancel()
	// getHistoricalPrices returns the results of historical price requests, which are possible to be collected
	// at this point in time
	historicalPrices, err := getHistoricalPrices(cctx)
	if err != nil {
		a.Logger.Error("getting the historical price failed: "+err.Error(), zap.String("requestID", requestID))
		return nil, err
	}

	investment := domain.Investment{
		FiatName: fiatName,
		Amount:   input.Amount,
	}
	lostChance, err := domain.CalculateLostChance(investment, historicalPrices, currentPrices)
	if err != nil {
		return nil, err
	}

	return &lostChance, nil
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
