package app

import (
	"context"
	"lost-chances-calc/internal/domain"
	"time"
)

type getHistoricalPricesFunc func(ctx context.Context) (prices []domain.HistoricalPrice, err error)

func (a App) requestHistoricalPrices(ctx context.Context, requestID string, progress *int, fiatName string, monthYear time.Time) (getterFunc getHistoricalPricesFunc, err error) {

	for _, currency := range domain.Cryptocurrencies {
		// call to fetch dispaches a task and subscribes to the result queue
		a.priceFetcher.FetchHistoricalPrice(ctx, requestID, currency, fiatName, monthYear)

		a.publishProgress(ctx, progress, requestID)
	}

	getterFunc = func(ctx context.Context) (prices []domain.HistoricalPrice, err error) {
		prices, err = a.priceFetcher.CollectHistoricalPrices(ctx, requestID, len(domain.Cryptocurrencies))

		if err != nil {
			return
		}
		a.publishProgress(ctx, progress, requestID)

		return
	}

	return getterFunc, nil

}
