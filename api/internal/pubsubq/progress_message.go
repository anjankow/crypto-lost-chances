package pubsubq

import (
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type ProgressMessage struct {
	RequestID string `json:"requestID"`
	Progress  int    `json:"progress"`
}

func GetProgressMessage(msg *pubsub.Message) (ProgressMessage, error) {
	var progressMessage ProgressMessage
	err := json.Unmarshal(msg.Data, &progressMessage)
	return progressMessage, err
}
