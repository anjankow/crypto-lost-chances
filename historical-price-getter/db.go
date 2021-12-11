// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/api/iterator"
)

const (
	dbName      = "db"
	dbInstance  = "internal"
	projectName = "crypto-lost-chances"
)

func getFromDB(ctx context.Context, h HistoricalPrice) (found bool, out HistoricalPrice) {
	found = false

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
		FROM prices WHERE
		cryptocurrency = @cryptocurrency
		AND
		fiat = @fiat
		AND
		monthYear = @monthYear
	`
	stmt := spanner.Statement{
		SQL: selectQ,
		Params: map[string]interface{}{
			"cryptocurrency": h.CryptocurrencyName,
			"fiat":           h.FiatName,
			"monthYear":      h.MonthYear,
		},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	// we need just one record
	row, err := iter.Next()
	if err == iterator.Done {
		fmt.Println("price doesn't exist in the db yet")
		return
	}
	if err != nil {
		fmt.Println("error when getting next result: ", err.Error())
		return
	}

	fmt.Println("price found in the db")

	if err := row.Columns(&out.CryptocurrencyName, &out.FiatName, &out.MonthYear,
		&out.PriceHighest, &out.PriceLowest); err != nil {
		fmt.Println("scanning error: ", err.Error())
	}
	fmt.Println("price: ", out.CryptocurrencyName, out.FiatName, out.MonthYear, out.PriceHighest, out.PriceLowest)

	return true, out

}
