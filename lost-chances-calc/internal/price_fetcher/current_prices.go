package pricefetcher

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"lost-chances-calc/internal/domain"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

const (
	providerURL = "https://api.bitpanda.com/v1/ticker"
)

type currentFiatPrices map[string]string

func (p *PriceFetcher) requestCurrentPrices() (raw json.RawMessage, err error) {

	resp, err := http.Get(providerURL)
	if err != nil {
		err = errors.New("failed to get current prices: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("error status code received from the provider: " + resp.Status)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("error while reading the provider's response: " + err.Error())
		return
	}

	if jsonErr := json.Unmarshal(b, &raw); err != nil {
		err = errors.New("error while unmarshalling the provider's response to raw json: " + jsonErr.Error())
		return
	}
	return

}

func (p *PriceFetcher) FetchCurrentPrices(ctx context.Context, requestID string) (domain.CurrentPricesDict, error) {
	raw, err := p.requestCurrentPrices()
	if err != nil {
		return nil, err
	}

	rawPrices := map[string]currentFiatPrices{}
	if err := json.Unmarshal(raw, &rawPrices); err != nil {
		return nil, errors.New("failed to unmarshal to raw prices: " + err.Error())
	}

	result := make(domain.CurrentPricesDict, len(domain.Fiats))
	for _, fiat := range domain.Fiats {
		result[fiat] = map[domain.Cryptocurrency]float64{}
	}

	for _, currency := range domain.Cryptocurrencies {
		pricesPerCurrency, ok := rawPrices[string(currency)]
		if !ok {
			p.logger.Warn("can't find the current prices of the cryptocurrency " + string(currency))
			continue
		}

		for _, fiat := range domain.Fiats {
			price, ok := pricesPerCurrency[string(fiat)]
			if !ok {
				p.logger.Warn("can't find the current " + string(fiat) + " price of the currency " + string(currency))
				continue
			}

			fPrice, err := strconv.ParseFloat(price, 64)
			if err != nil {
				p.logger.Warn("can't parse current price to float: "+err.Error(), zap.String("current_price", price))
			}
			result[fiat][currency] = fPrice
		}

	}

	return result, nil
}
