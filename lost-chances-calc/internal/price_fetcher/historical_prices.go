package pricefetcher

import (
	"context"
	"lost-chances-calc/internal/domain"
	"time"

	"go.uber.org/zap"
)

const (
	waitTime = 200 * time.Millisecond
)

func (p *PriceFetcher) FetchHistoricalPrice(ctx context.Context, requestID string, cryptoCurr domain.Cryptocurrency, fiatName string, monthYear time.Time) {
	// submits a task to fetch the data

	// subscribes to the results
	p.subscribeToHistoricalPrices(requestID)

}

func (p *PriceFetcher) CollectHistoricalPrices(ctx context.Context, requestID string, expectedNum int) ([]domain.HistoricalPrice, error) {
	var prices []domain.HistoricalPrice

	for {
		prices = p.getHistoricalPrices(requestID)
		if len(prices) == expectedNum {
			break
		}
		if ctx.Err() != nil {
			p.logger.Warn("timeout when collecting the historical prices", zap.String("requestID", requestID))
			break
		}
		time.Sleep(waitTime)
	}

	p.unsubscribeFromHistoricalPrices(requestID)

	return prices, nil
}
