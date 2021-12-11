package app

import (
	"context"
	"lost-chances-calc/internal/domain"
	"time"
)

func (a App) getHistoricalPrices(ctx context.Context, requestID string, progress *int, fiatName string, monthYear time.Time) (prices []domain.HistoricalPrice, err error) {

	for _, currency := range domain.Cryptocurrencies {
		a.priceFetcher.FetchHistoricalPrice(ctx, requestID, currency, fiatName, monthYear)

		*progress += progressStepLen
		a.progressWriter.PublishProgress(ctx, requestID, *progress)
	}

	cctx, cancel := context.WithTimeout(ctx, priceFetcherResultsTimeout)
	defer cancel()

	for range domain.Cryptocurrencies {
		price, err := a.priceFetchListener.GetPrice(cctx, requestID)
		if err != nil {

			return prices, err
		}

		*progress += progressStepLen
		a.progressWriter.PublishProgress(ctx, requestID, *progress)

		prices = append(prices, price)
	}

	return
}
