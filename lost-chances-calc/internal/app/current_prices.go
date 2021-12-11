package app

import (
	"context"
	"lost-chances-calc/internal/domain"
)

func (a App) getCurrentPrices(ctx context.Context, fiatName domain.Fiat) (prices []domain.CurrentPrice, err error) {
	return
}
