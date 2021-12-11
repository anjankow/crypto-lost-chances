package p

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	providerURL = "https://min-api.cryptocompare.com/data/v2/histoday"
)

type singlePriceData struct {
	High float64 `json:"high"`
	Low  float64 `json:"low"`
}

type priceData struct {
	Data []singlePriceData `json:"Data"`
}

type providerRsp struct {
	ResponseStatus string    `json:"Response"`
	Data           priceData `json:"Data"`
}

func checkPrice(input HistoricalPrice) (output HistoricalPrice, err error) {

	timeFrom := time.Date(input.MonthYear.Year(), input.MonthYear.Month(), 1, 0, 0, 0, 0, time.Local)
	timeTo := timeFrom.AddDate(0, 1, -1)

	req, err := http.NewRequest(http.MethodGet, providerURL, nil)
	if err != nil {
		err = errors.New("failed to create a new request: " + err.Error())
		return
	}

	q := req.URL.Query()
	q.Add("fsym", input.CryptocurrencyName)
	q.Add("tsym", input.FiatName)

	q.Add("limit", "1")
	q.Add("aggregate", "30")
	q.Add("toTs", fmt.Sprint(timeTo.Unix()))

	req.URL.RawQuery = q.Encode()
	log.Printf("requesting historical data: " + req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New("failed to request historical data: " + err.Error())
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

	presp := providerRsp{}
	if jsonErr := json.Unmarshal(b, &presp); err != nil {
		err = errors.New("error while unmarshalling the provider's response: " + jsonErr.Error())
		return
	}

	if presp.ResponseStatus != "Success" {
		err = errors.New("provider's response status code is not `Success`: " + presp.ResponseStatus)
		return
	}

	if len(presp.Data.Data) < 1 {
		err = errors.New("no data in the provider's response")
		return
	}

	output = HistoricalPrice{
		CryptocurrencyName: input.CryptocurrencyName,
		FiatName:           input.FiatName,
		MonthYear:          input.MonthYear,
	}
	output.PriceHighest = presp.Data.Data[0].High
	output.PriceLowest = presp.Data.Data[0].Low

	return
}
