package p

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/multierr"
)

var (
	// arbitrary limit - one month after ETH was released
	pastLimit = time.Date(2015, time.August, 1, 0, 0, 0, 0, time.Local)
)

type HistoricalPrice struct {
	CryptocurrencyName string    `json:"cryptocurrency"`
	FiatName           string    `json:"fiat"`
	MonthYear          time.Time `json:"monthYear"`
	PriceHighest       float64   `json:"priceHighest"`
	PriceLowest        float64   `json:"priceLowest"`
}

func (h HistoricalPrice) Validate() (err error) {

	if h.CryptocurrencyName == "" {
		err = multierr.Append(err, errors.New("missing cryptocurrency name"))
	}
	if h.FiatName == "" {
		err = multierr.Append(err, errors.New("missing fiat name"))
	}
	err = multierr.Append(err, h.validateMonthYear())

	return err
}

func (h HistoricalPrice) validateMonthYear() error {
	if h.MonthYear.IsZero() {
		return errors.New("missing date")
	}
	if h.MonthYear.Before(pastLimit) {
		return fmt.Errorf("the date is too much in the past, limit is %v", pastLimit)
	}
	if h.MonthYear.After(time.Now()) {
		return errors.New("the date has to be from the past")
	}

	return nil
}
