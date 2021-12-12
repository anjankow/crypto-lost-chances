package domain

import "time"

type HistoricalPrice struct {
	CryptocurrencyName Cryptocurrency `json:"cryptocurrency"`
	FiatName           Fiat           `json:"fiat"`
	MonthYear          time.Time      `json:"monthYear"`
	PriceHighest       float64        `json:"priceHighest"`
	PriceLowest        float64        `json:"priceLowest"`
}
