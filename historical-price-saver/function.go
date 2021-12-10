// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	user                   = "root"
	password               = ""
	instanceConnectionName = "crypto-lost-chances:europe-central2:db"
	dbName                 = "test"
	socketDir              = "/cloudsql"
)

type PubSubMessage struct {
	RequestID string          `json:"requestID"`
	Price     HistoricalPrice `json:"historicalPrice"`
}

type HistoricalPrice struct {
	CryptocurrencyName string    `json:"cryptocurrency"`
	FiatName           string    `json:"fiat"`
	MonthYear          time.Time `json:"monthYear"`
	PriceHighest       float64   `json:"priceHighest"`
	PriceLowest        float64   `json:"priceLowest"`
}

func SavePrice(ctx context.Context, m PubSubMessage) error {

	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", user, password, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		fmt.Println("error in db connection: " + err.Error())
		return err
	}

	selectQ := `
	SELECT * FROM prices WHERE
		cryptocurrency = ?
		AND
		fiat = ?
		AND
		monthYear = ?;
	`

	h := m.Price
	row := dbPool.QueryRowContext(ctx, selectQ, h.FiatName, h.CryptocurrencyName, h.MonthYear)
	result := HistoricalPrice{}
	if err := row.Scan(&result); err != sql.ErrNoRows {
		if err == nil {
			fmt.Println("price already exists in the db: ", h.CryptocurrencyName, "/", h.FiatName, " ", h.MonthYear)
			return nil
		}
		fmt.Println("error when querying the db: " + err.Error())
		return err
	}

	query := `
	INSERT INTO prices(
		cryptocurrency,
		fiat,
		monthYear,
		priceHighest,
		priceLowest
	) VALUES (
		?, ?, ?, ?, ?
	)`

	_, err = dbPool.ExecContext(ctx, query, h.FiatName, h.CryptocurrencyName, h.MonthYear)
	if err != nil {
		fmt.Println("error when inserting a price: " + err.Error())
		return err
	}

	return nil
}
