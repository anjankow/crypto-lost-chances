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
)

var (
	url         = config.GetLostChancesCalcHost() + "/calculate"
	contentType = "application/json"
)

type calcRequestBody struct {
	RequestID string    `json:"requestID"`
	MonthYear time.Time `json:"monthYear`
	Amount    int       `json:"amount"`
}

type Client struct {
	logger *zap.Logger
}

func NewClient(logger *zap.Logger) Client {
	return Client{logger: logger}
}

func (c Client) StartCalculation(ctx context.Context, requestID string, monthYear time.Time, amount int) error {

	body := calcRequestBody{
		RequestID: requestID,
		MonthYear: monthYear,
		Amount:    amount,
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return errors.New("failed to marshal the body: " + err.Error())
	}
	resp, err := http.Post(url, contentType, bytes.NewBuffer(marshalledBody))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status is " + resp.Status)
	} else {
		c.logger.Debug("calculation request succeeded", zap.Any("status_code", resp.StatusCode), zap.String("requestID", requestID))
	}

	return nil
}
