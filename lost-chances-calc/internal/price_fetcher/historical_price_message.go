package pricefetcher

import (
	"encoding/json"
	"errors"
	"lost-chances-calc/internal/domain"
)

type HistoricalPriceMessage struct {
	RequestID string                 `json:"requestID"`
	Price     domain.HistoricalPrice `json:"historicalPrice"`
}

func unmarshalHistoricalPriceMessage(message []byte) (HistoricalPriceMessage, error) {

	var h HistoricalPriceMessage
	if err := json.Unmarshal(message, &h); err != nil {
		return h, errors.New("failed to unmarshal historical price message: " + err.Error())
	}

	return h, nil
}
