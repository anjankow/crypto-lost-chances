package pubsubq

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Config struct {
	Topic  string
	SubID  string
	Filter string
}

func Subscribe(ctx context.Context, client pubsub.Client, cfg Config) error {

	subCfg := pubsub.SubscriptionConfig{
		Topic:                 client.Topic(cfg.Topic),
		EnableMessageOrdering: true,
		Filter:                cfg.Filter,
		Detached:              false,
	}

	sub, err := client.CreateSubscription(ctx, cfg.SubID, subCfg)
	if err != nil {
		return err
	}

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	return nil
}
