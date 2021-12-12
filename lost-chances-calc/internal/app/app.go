package app

import (
	"context"
	"errors"
	"lost-chances-calc/internal/domain"
	pricefetcher "lost-chances-calc/internal/price_fetcher"
	progressupdates "lost-chances-calc/internal/progress_updates"
	"time"

	"go.uber.org/zap"
)

const (
	priceFetcherResultsTimeout = 5 * time.Second
	progressMax                = 100
	// support only euro for now
	supportedFiatName = domain.Euro
)

var (
	progressSteps = len(domain.Cryptocurrencies) + //for submitting fetcher tasks
		1 + // for getting the historical prices
		1 + // for getting the current price
		1 + // for the initial progress added just because the operation reached this point
		1 // for calculating the final result

	progressStepLen = int(progressMax / progressSteps)
)

type App struct {
	Logger         *zap.Logger
	progressWriter *progressupdates.Writer
	priceFetcher   *pricefetcher.PriceFetcher
}

type CalcInput struct {
	MonthYear time.Time
	Amount    float64
}

type Results struct {
	Cryptocurrency string `json:"currency"`
	Income         string `json:"income"`
}

func NewApp(l *zap.Logger, progressWriter *progressupdates.Writer, fetcher *pricefetcher.PriceFetcher) (app App, err error) {
	if progressWriter == nil {
		err = errors.New("progress writer is nil")
		return
	}

	app = App{
		Logger:         l,
		progressWriter: progressWriter,
		priceFetcher:   fetcher,
	}
	return
}

func (a App) Calculate(ctx context.Context, requestID string, input CalcInput) (chance *domain.LostChance, err error) {

	progress := 0
	// if we are here, then some work has been done already - update the progress
	a.publishProgress(ctx, &progress, requestID)

	// request to dispatch the getter tasks and subscribe for the historical price messages
	getHistoricalPrices, err := a.requestHistoricalPrices(ctx, requestID, &progress, supportedFiatName, input.MonthYear)
	if err != nil {
		a.Logger.Error("requesting the historical prices failed: "+err.Error(), zap.String("requestID", requestID))
	}

	currentPrices, err := a.getCurrentPrices(ctx, requestID, &progress)
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
	a.publishProgress(ctx, &progress, requestID)

	investment := domain.Investment{
		FiatName: supportedFiatName,
		Amount:   input.Amount,
	}
	lostChance, err := domain.CalculateLostChance(investment, historicalPrices, currentPrices)
	if err != nil {
		return nil, err
	}
	a.publishOperationCompleted(requestID)

	return &lostChance, nil
}

func (a App) publishProgress(ctx context.Context, currentProgress *int, requestID string) {

	if *currentProgress+progressStepLen < progressMax {
		*currentProgress += progressStepLen
	}

	a.progressWriter.PublishProgress(ctx, requestID, *currentProgress)
}

func (a App) publishOperationCompleted(requestID string) {
	a.progressWriter.PublishProgress(context.Background(), requestID, progressMax)
}
