package p

import (
	"context"
	"encoding/json"
	"errors"

	"cloud.google.com/go/pubsub"
)

const (
	projectID = "crypto-lost-chances"

	topicName = "historical-prices"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func writeResult(ctx context.Context, h HistoricalPriceMessage) error {

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return errors.New("failed to create a pubsub client: " + err.Error())
	}
	defer client.Close()

	topic := client.Topic(topicName)
	if topic == nil {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return errors.New("failed to create a new topic: " + err.Error())
		}
	}

	bytes, err := json.Marshal(h)
	if err != nil {
		return errors.New("error when marshalling the message: " + err.Error())
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: bytes,
	})

	if _, err := result.Get(ctx); err != nil {
		return errors.New("error when publishing the message: " + err.Error())
	}

	return nil
}
