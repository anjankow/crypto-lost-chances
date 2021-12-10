// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type requestBody struct {
	RequestID    string          `json:"requestID"`
	PriceRequest HistoricalPrice `json:"priceRequest"`
}

func GetHistoricalPrice(w http.ResponseWriter, r *http.Request) {

	var input requestBody
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

	if input.RequestID == "" {
		http.Error(w, "missing request ID", http.StatusBadRequest)
		return
	}

	if err := input.PriceRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	found, price := getFromDB(r.Context(), input.PriceRequest)
	if !found {
		h, err := checkPrice(input.PriceRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		price = h
	}

	message := fmt.Sprintln(price.CryptocurrencyName, "/", price.FiatName, " lowest: ", price.PriceLowest, ", highest: ", price.PriceHighest)
	fmt.Fprint(w, message)
}
