package pubsubq

import (
	"context"

	"cloud.google.com/go/pubsub"
)

const (
	subscriptionName = "progressSub"
	topicName        = "progress-update"
)

func Subscribe(ctx context.Context, client *pubsub.Client) (*pubsub.Subscription, error) {

	sub := client.Subscription(subscriptionName)
	if sub != nil {
		return sub, nil
	}

	subCfg := pubsub.SubscriptionConfig{
		Topic:                 client.Topic(topicName),
		EnableMessageOrdering: true,
		Detached:              false,
	}

	sub, err := client.CreateSubscription(ctx, subscriptionName, subCfg)
	if err != nil {
		return nil, err
	}

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	return sub, nil
}
