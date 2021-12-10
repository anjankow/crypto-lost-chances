package p

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	providerURL = "https://min-api.cryptocompare.com/data/v2/histoday"
)

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
	}

	resp.Body.Close()
	return
}
