// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"
)

const (
	user                   = "root"
	password               = ""
	instanceConnectionName = "crypto-lost-chances:europe-central2:db"
	dbName                 = "internal"
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

func (h HistoricalPrice) Validate() (err error) {

	if h.CryptocurrencyName == "" {
		err = multierr.Append(err, errors.New("missing cryptocurrency name"))
	}
	if h.FiatName == "" {
		err = multierr.Append(err, errors.New("missing fiat name"))
	}
	if h.MonthYear.IsZero() {
		err = multierr.Append(err, errors.New("missing date"))
	}
	if h.PriceHighest == 0 && h.PriceLowest == 0 {
		err = multierr.Append(err, errors.New("missing prices lowest and highest"))
	}

	return err
}

func SavePrice(ctx context.Context, m PubSubMessage) error {

	h := m.Price
	if err := h.Validate(); err != nil {
		fmt.Println("validation failed: " + err.Error())
		return err
	}

	fmt.Println(m.RequestID, m.Price, m.Price.MonthYear)
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

	_, err = dbPool.ExecContext(ctx, query, h.FiatName, h.CryptocurrencyName, h.MonthYear, h.PriceHighest, h.PriceLowest)
	if err != nil {
		fmt.Println("error when inserting a price: " + err.Error())
		return err
	}

	return nil
}
