// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"google.golang.org/api/iterator"
)

const (
	dbName      = "db"
	dbInstance  = "internal"
	projectName = "crypto-lost-chances"
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
	if err := json.Unmarshal(message.Data, &m); err != nil {
		return errors.New("unmarshalling the message failed: " + err.Error())
	}

	fmt.Println("saving, request id: ", m.RequestID, ", price: ", m.Price)

	h := m.Price
	if err := h.Validate(); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	dsn := fmt.Sprint("projects/", projectName, "/instances/", dbInstance, "/databases/", dbName)
	client, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		return err
	}
	defer client.Close()

	selectQ := `
	SELECT 
		cryptocurrency,
		fiat,
		monthYear
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
	for {
		row, err := iter.Next()
		fmt.Println("got next iterator")
		if err == iterator.Done {
			fmt.Println("iterator done")
			break
		}
		if err != nil {
			fmt.Println("iterator done")
			return err
		}
		var histRead HistoricalPrice
		if err := row.Columns(nil, &histRead.CryptocurrencyName, &histRead.FiatName, &histRead.MonthYear); err != nil {
			fmt.Println("scanning error: ", err.Error())
		}
		fmt.Println("hist price: ", histRead.CryptocurrencyName, histRead.FiatName, histRead.MonthYear)
	}

	fmt.Println("now inserting")

	query := `
	INSERT INTO prices(
		id,
		cryptocurrency,
		fiat,
		monthYear,
		priceHighest,
		priceLowest
	) VALUES (
		@uuid, @cryptocurrency, @fiat, @monthYear, @priceHighest, @priceLowest
	)`
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: query,
			Params: map[string]interface{}{
				"uuid":           uuid.New().String(),
				"cryptocurrency": h.CryptocurrencyName,
				"fiat":           h.FiatName,
				"monthYear":      h.MonthYear,
				"priceHighest":   h.PriceHighest,
				"priceLowest":    h.PriceLowest,
			},
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		if rowCount != 1 {
			return errors.New("one record should be inserted, executed row count:" + fmt.Sprint(rowCount))
		}

		fmt.Println("saved, request id: ", m.RequestID, ", price: ", m.Price)
		return err
	})

	fmt.Println("return")

	return err
}
