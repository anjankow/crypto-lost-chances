// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HistoricalPriceMessage struct {
	RequestID string          `json:"requestID"`
	Price     HistoricalPrice `json:"historicalPrice"`
}

func GetHistoricalPrice(w http.ResponseWriter, r *http.Request) {

	var input HistoricalPriceMessage
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

	if err := input.Price.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	found, price := getFromDB(r.Context(), input.Price)
	if !found {
		h, err := checkPrice(input.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		price = h
	}

	fmt.Println(price.CryptocurrencyName, "/", price.FiatName, " lowest: ", price.PriceLowest, ", highest: ", price.PriceHighest)

	writeResult(r.Context(), HistoricalPriceMessage{RequestID: input.RequestID, Price: price})
}
