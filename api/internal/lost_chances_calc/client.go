package lostchancescalc

import (
	"api/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
)

var (
	url         = config.GetLostChancesCalcHost() + "/calculate"
	contentType = "application/json"
)

type calcRequestBody struct {
	RequestID string    `json:"requestID"`
	MonthYear time.Time `json:"monthYear"`
	Amount    int       `json:"amount"`
}

type Client struct {
	logger *zap.Logger
}

func NewClient(logger *zap.Logger) Client {
	return Client{logger: logger}
}

func (c Client) Calculate(ctx context.Context, requestID string, monthYear time.Time, amount int) (lostChance LostChance, err error) {

	body := calcRequestBody{
		RequestID: requestID,
		MonthYear: monthYear,
		Amount:    amount,
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		err = errors.New("failed to marshal the body: " + err.Error())
		return
	}

	client, err := idtoken.NewClient(ctx, url)
	if err != nil {
		err = errors.New("can't create a idtoken client: " + err.Error())
		return
	}

	c.logger.Debug("posting calculation request to "+url, zap.String("requestID", requestID))
	resp, err := client.Post(url, contentType, bytes.NewBuffer(marshalledBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("response status is " + resp.Status)
		return
	}
	c.logger.Debug("calculation request succeeded", zap.Any("status_code", resp.StatusCode), zap.String("requestID", requestID))

	if err = json.NewDecoder(resp.Body).Decode(&lostChance); err != nil {
		err = errors.New("failed to decode response: " + err.Error())
		return
	}

	c.logger.Debug("lost chance received", zap.Any("lost_chance", lostChance), zap.String("requestID", requestID))

	return
}
