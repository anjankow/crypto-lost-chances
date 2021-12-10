package p

import (
	"errors"
	"time"

	"go.uber.org/multierr"
)

type HistoricalPrice struct {
	RequestID          string    `json:"requestID"`
	CryptocurrencyName string    `json:"cryptocurrency"`
	FiatName           string    `json:"fiatName"`
	MonthYear          time.Time `json:"monthYear"`
	PriceHighest       float64   `json:"priceHighest"`
	PriceLowest        float64   `json:"priceLowest"`
}

func (h HistoricalPrice) Validate() (err error) {
	if h.RequestID == "" {
		err = multierr.Append(err, errors.New("missing request ID"))
	}
	if h.CryptocurrencyName == "" {
		err = multierr.Append(err, errors.New("missing cryptocurrency name"))
	}
	if h.FiatName == "" {
		err = multierr.Append(err, errors.New("missing fiat name"))
	}
	if h.MonthYear.IsZero() {
		err = multierr.Append(err, errors.New("missing date"))
	}

	return err
}
