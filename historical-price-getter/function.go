// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetHistoricalPrice(w http.ResponseWriter, r *http.Request) {

	var input HistoricalPrice
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		switch err {
		case io.EOF:
			fmt.Fprint(w, "empty body")
			return
		default:
			log.Printf("decoder error: %v", err)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	historicalPrice, err := checkPrice(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message := fmt.Sprintln(historicalPrice.CryptocurrencyName, "/", historicalPrice.FiatName, " lowest: ", historicalPrice.PriceLowest, ", highest: ", historicalPrice.PriceHighest)
	fmt.Fprint(w, message)
}
