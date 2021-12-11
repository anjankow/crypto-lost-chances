// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/spanner"
	"go.uber.org/multierr"
)

const (
	dbName      = "db"
	dbInstance  = "internal"
	projectName = "crypto-lost-chances"
)

type requestInfo struct {
	RequestID string    `json:"requestID"`
	FiatName  string    `json:"fiat"`
	MonthYear time.Time `json:"monthYear"`
}

func (r requestInfo) validate() (err error) {

	if r.RequestID == "" {
		err = multierr.Append(err, errors.New("missing request ID"))
	}
	if r.FiatName == "" {
		err = multierr.Append(err, errors.New("missing fiat name"))
	}
	if r.MonthYear.IsZero() {
		err = multierr.Append(err, errors.New("missing date"))
	}

	return err
}

func saveInDB(ctx context.Context, r requestInfo) error {
	dsn := fmt.Sprint("projects/", projectName, "/instances/", dbInstance, "/databases/", dbName)
	client, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		return errors.New("error when connecting to the db: " + err.Error())
	}
	defer client.Close()

	query := `
	INSERT INTO requests(
		id,
		fiat,
		monthYear
	) VALUES (
		@requestID, @fiat, @monthYear
	)`

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: query,
			Params: map[string]interface{}{
				"id":        r.RequestID,
				"fiat":      r.FiatName,
				"monthYear": r.MonthYear,
			},
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		if rowCount != 1 {
			return errors.New("one record should be inserted, executed row count:" + fmt.Sprint(rowCount))
		}

		fmt.Println("saved, request id: ", r.RequestID)
		return err
	})

	return err
}

func SaveRequest(w http.ResponseWriter, r *http.Request) {

	var input requestInfo
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

	if err := input.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := saveInDB(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "request saved successfully, request ID: "+input.RequestID)
}
