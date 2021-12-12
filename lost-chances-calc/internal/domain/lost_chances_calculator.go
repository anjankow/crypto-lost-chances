package domain

import (
	"errors"
	"fmt"
)

func CalculateLostChance(fiatInvestment Investment, historicalPrices []HistoricalPrice, currentPrices CurrentPricesDict) (chance LostChance, err error) {

	chance.FiatName = fiatInvestment.FiatName
	chance.InputFiatAmount = fiatInvestment.Amount

	currentPricesForFiat, ok := currentPrices[fiatInvestment.FiatName]
	if !ok {
		err = errors.New("no current crytpo prices for fiat " + string(fiatInvestment.FiatName))
		return
	}

	outAmounts := map[Cryptocurrency]float64{}
	for _, h := range historicalPrices {
		if h.FiatName != fiatInvestment.FiatName {
			// fiats don't match, can't compare them
			continue
		}

		currentPrice, ok := currentPricesForFiat[h.CryptocurrencyName]
		if !ok {
			continue
		}

		outAmounts[h.CryptocurrencyName] = getFiatOutAmount(fiatInvestment.Amount, h.PriceHighest, currentPrice)

	}

	if len(outAmounts) == 0 {
		err = fmt.Errorf("can't calculate any lost chances, len(historicalPrices): %v, len(currentPrices): %v", len(historicalPrices), len(currentPrices))
		return
	}

	for key, val := range outAmounts {
		if val > chance.OutputFiatAmount {
			chance.CryptocurName = key
			chance.OutputFiatAmount = val
		}
	}

	return
}

func getFiatOutAmount(fiatInAmount float64, historicalPrice float64, currentPrice float64) float64 {
	cryptoAmount := fiatInAmount / historicalPrice

	return currentPrice * cryptoAmount
}
