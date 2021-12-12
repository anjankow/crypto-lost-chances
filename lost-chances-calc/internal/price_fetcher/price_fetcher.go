package pricefetcher

import (
	"context"
	"lost-chances-calc/internal/domain"
	"time"

	"go.uber.org/zap"
)

type PriceFetcher struct {
	logger *zap.Logger
}

func NewFetcher(logger *zap.Logger) PriceFetcher {
	return PriceFetcher{logger: logger}
}

func (p *PriceFetcher) Init(ctx context.Context) (closer func(), err error) {
	return
}

func (p *PriceFetcher) FetchHistoricalPrice(ctx context.Context, requestID string, cryptoCurr domain.Cryptocurrency, fiatName string, monthYear time.Time) {
	// submits a task to fetch the data
}

func (p *PriceFetcher) CollectHistoricalPrices(ctx context.Context, requestID string) ([]domain.HistoricalPrice, error) {
	return []domain.HistoricalPrice{}, nil
}

func (p *PriceFetcher) FetchCurrentPrices(ctx context.Context, requestID string) (domain.CurrentPricesDict, error) {
	return domain.CurrentPricesDict{}, nil
}
