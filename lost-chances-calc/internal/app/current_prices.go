package app

import (
	"context"
	"lost-chances-calc/internal/domain"
)

func (a App) getCurrentPrices(ctx context.Context, requestID string, progress *int) (prices domain.CurrentPricesDict, err error) {

	if prices, err = a.priceFetcher.FetchCurrentPrices(ctx, requestID); err != nil {
		return nil, err
	}

	a.publishProgress(ctx, progress, requestID)
	return prices, err
}
