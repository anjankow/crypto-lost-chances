// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	user                   = "root"
	password               = ""
	instanceConnectionName = "crypto-lost-chances:europe-central2:db"
	dbName                 = "test"
	socketDir              = "/cloudsql"
)

func getFromDB(ctx context.Context, h HistoricalPrice) (found bool, out HistoricalPrice) {

	found = false
	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", user, password, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		fmt.Println("error in db connection: " + err.Error())
		return
	}

	query := `
	SELECT * FROM prices 
	WHERE
	fiat = ?
	AND
	cryptocurrency = ?
	AND
	date = ?
	`
	rows, err := dbPool.QueryContext(ctx, query, h.FiatName, h.CryptocurrencyName, h.MonthYear)
	if err != nil {
		fmt.Println("error when executing the query: " + err.Error())
		return
	}

	var prices []HistoricalPrice
	err = rows.Scan(prices)
	if err != nil {
		fmt.Println("row scan failed: " + err.Error())
		return
	}

	if len(prices) == 0 {
		fmt.Println("no prices found in the db")
		return
	}

	return true, prices[0]
}
