package dbviewer

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
)

const (
	dbViewerURL = "https://europe-central2-crypto-lost-chances.cloudfunctions.net/db-viewer"
)

func GetAllHistoricalPrices(ctx context.Context, logger *zap.Logger) (rawRecords []byte, err error) {

	client, err := idtoken.NewClient(ctx, dbViewerURL)
	if err != nil {
		err = errors.New("can't create a idtoken client: " + err.Error())
	}

	logger.Debug("getting the db records")
	resp, err := client.Get(dbViewerURL)
	if err != nil {
		err = errors.New("can't access db viewer: " + err.Error())
		return
	}
	defer resp.Body.Close()

	logger.Debug("response received")
	if resp.StatusCode != http.StatusOK {
		err = errors.New("response status is " + resp.Status)
		return
	}

	rawRecords, err = ioutil.ReadAll(resp.Body)
	if len(rawRecords) == 0 {
		err = errors.New("received an empty response")
	}

	return
}
