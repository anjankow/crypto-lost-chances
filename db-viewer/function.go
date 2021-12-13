// Package p contains an HTTP Cloud Function.
package p

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const (
	dbName      = "db"
	dbInstance  = "internal"
	projectName = "crypto-lost-chances"
)

type HistoricalPrice struct {
	CryptocurrencyName string    `json:"cryptocurrency"`
	FiatName           string    `json:"fiat"`
	MonthYear          time.Time `json:"monthYear"`
	PriceHighest       float64   `json:"priceHighest"`
	PriceLowest        float64   `json:"priceLowest"`
}

func ViewDB(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	dsn := fmt.Sprint("projects/", projectName, "/instances/", dbInstance, "/databases/", dbName)
	client, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		fmt.Println("error when connecting to the db: ", err.Error())
		return
	}
	defer client.Close()

	selectQ := `
	SELECT
		cryptocurrency,
		fiat,
		monthYear,
		priceHighest,
		priceLowest
		FROM prices
	`
	stmt := spanner.Statement{
		SQL: selectQ,
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var records []HistoricalPrice

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			// end of records
			break
		}
		if err != nil {
			fmt.Println("error when getting a next result: ", err.Error())
			continue
		}

		var out HistoricalPrice
		if err := row.Columns(&out.CryptocurrencyName, &out.FiatName, &out.MonthYear,
			&out.PriceHighest, &out.PriceLowest); err != nil {
			fmt.Println("scanning error: ", err.Error())
			continue
		}

		fmt.Println("price: ", out.CryptocurrencyName, out.FiatName, out.MonthYear, out.PriceHighest, out.PriceLowest)
		records = append(records, out)
	}

	if len(records) == 0 {
		fmt.Println("no records in the db")
		fmt.Fprintf(w, "no records in the db")
		return
	}

	b, err := json.Marshal(records)
	if err != nil {
		fmt.Println("failed to marshal the records: " + err.Error())
		return
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write the response: " + err.Error())
	}
}
