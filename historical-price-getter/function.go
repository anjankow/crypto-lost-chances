// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"html"
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
			log.Printf("json.NewDecoder: %v", err)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprint(w, html.EscapeString(input.CryptocurrencyName))
}
