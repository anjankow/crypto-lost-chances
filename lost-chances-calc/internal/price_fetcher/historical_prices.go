package pricefetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lost-chances-calc/internal/domain"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"go.uber.org/zap"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

const (
	waitTime = 200 * time.Millisecond

	queueID    = "historical-price-req"
	projectID  = "crypto-lost-chances"
	locationID = "europe-central2"

	priceGetterURL = "https://europe-central2-crypto-lost-chances.cloudfunctions.net/historical-price-getter"
)

var (
	queuePath = fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)
)

func (p *PriceFetcher) FetchHistoricalPrice(ctx context.Context, requestID string, cryptoCurr domain.Cryptocurrency, fiatName domain.Fiat, monthYear time.Time) {
	// submits a task to fetch the data

	message := HistoricalPriceMessage{
		RequestID: requestID,
		Price: domain.HistoricalPrice{
			CryptocurrencyName: cryptoCurr,
			FiatName:           fiatName,
			MonthYear:          monthYear,
			PriceHighest:       0,
			PriceLowest:        0,
		},
	}
	if err := p.submitTask(ctx, message); err != nil {
		p.logger.Warn("failed to submit fetcher task: "+err.Error(), zap.String("requestID", requestID))
	}

	// subscribes to the results
	p.subscribeToHistoricalPrices(requestID)

}

func (p *PriceFetcher) submitTask(ctx context.Context, message HistoricalPriceMessage) error {

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	body, err := json.Marshal(message)
	if err != nil {
		return errors.New("failed to marshal historical price message: " + err.Error())
	}

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					Url:        priceGetterURL,
					HttpMethod: taskspb.HttpMethod_POST,
					Body:       body,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: "crypto-lost-chances@appspot.gserviceaccount.com",
						},
					},
				},
			},
		},
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return err
	}

	p.logger.Debug("submitted task", zap.Time("scheduled", createdTask.GetScheduleTime().AsTime()), zap.String("requestID", message.RequestID))

	return nil
}

func (p *PriceFetcher) CollectHistoricalPrices(ctx context.Context, requestID string, expectedNum int) ([]domain.HistoricalPrice, error) {
	var prices []domain.HistoricalPrice

	for {
		prices = p.getHistoricalPrices(requestID)
		if len(prices) == expectedNum {
			break
		}
		if ctx.Err() != nil {
			p.logger.Warn("timeout when collecting the historical prices", zap.String("requestID", requestID), zap.Int("prices_len", len(prices)))
			break
		}
		time.Sleep(waitTime)
	}

	p.unsubscribeFromHistoricalPrices(requestID)

	return prices, nil
}
