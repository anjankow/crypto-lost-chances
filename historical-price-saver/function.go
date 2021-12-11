// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"database/sql"
	"encoding/json"
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
	Data []byte `json:"data"`
}

type HistoricalPriceMessage struct {
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

func SavePrice(ctx context.Context, message PubSubMessage) error {

	var m HistoricalPriceMessage
	if err := json.Unmarshal(message.Data, &message); err != nil {
		return errors.New("unmarshalling the message failed: " + err.Error())
	}

	fmt.Println("saving, request id: ", m.RequestID, ", price: ", m.Price)

	h := m.Price
	if err := h.Validate(); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", user, password, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return errors.New("error in db connection: " + err.Error())
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

		return errors.New("error when querying the db: " + err.Error())
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

		return errors.New("error when inserting a price: " + err.Error())
	}

	return nil
}
