package pubsubq

import (
	"api/internal/config"
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
)

const (
	subscriptionName = "progressSub"
	topicName        = "progress-update"
)

func SubscribeToProgressUpdates(ctx context.Context) (sub *pubsub.Subscription, closerFunc func(), err error) {

	cctx, cancel := context.WithCancel(ctx)
	var client *pubsub.Client

	closerFunc = func() {
		cancel()
		client.Close()
	}

	client, err = pubsub.NewClient(cctx, config.GetProjectID())
	if err != nil {
		err = errors.New("failed to create a pubsub client: " + err.Error())
		return
	}

	sub = client.Subscription(subscriptionName)
	if sub != nil {
		// subscription already created, simply return
		return
	}

	// the subscription doesn't exist yet, create

	subCfg := pubsub.SubscriptionConfig{
		Topic:                 client.Topic(topicName),
		EnableMessageOrdering: true,
		Detached:              false,
	}

	sub, err = client.CreateSubscription(cctx, subscriptionName, subCfg)
	if err != nil {
		err = errors.New("failed to create a subscription: " + err.Error())
		return
	}

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	return
}
