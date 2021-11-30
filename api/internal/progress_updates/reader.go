package progressupdates

import (
	"api/internal/config"
	"api/internal/pubsubq"
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

const (
	subscriptionName = "progressSub"
	topicName        = "progress-update"

	maxProgress = 100
)

type Reader struct {
	logger         *zap.Logger
	sub            *pubsub.Subscription
	progressPerReq map[string]chan int
}

func NewReader(logger *zap.Logger) Reader {
	return Reader{
		logger:         logger,
		progressPerReq: map[string]chan int{},
	}
}

func (r *Reader) SubscribeToProgressUpdates(requestID string) chan int {
	_, ok := r.progressPerReq[requestID]
	if !ok {
		r.progressPerReq[requestID] = make(chan int, maxProgress)
	}

	return r.progressPerReq[requestID]
}

func (r *Reader) Start(ctx context.Context) (closer func(), err error) {

	cctx, ctxCancel := context.WithCancel(ctx)
	sub, subCloser, err := subscribe(cctx)
	closer = func() {
		subCloser()
		ctxCancel()
	}

	if err != nil {
		err = errors.New("subscribing to pubsub failed: " + err.Error())
		return
	}

	r.sub = sub

	go func() {
		if err := r.receiveFromPubsub(cctx); err != nil {
			r.logger.Error("error when receiving from pubsub: " + err.Error())
		}
		r.logger.Info("finished receiving from pubsub")
	}()

	return

}

func subscribe(ctx context.Context) (sub *pubsub.Subscription, closerFunc func(), err error) {

	var client *pubsub.Client

	closerFunc = func() {
		client.Close()
	}

	client, err = pubsub.NewClient(ctx, config.GetProjectID())
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

	sub, err = client.CreateSubscription(ctx, subscriptionName, subCfg)
	if err != nil {
		err = errors.New("failed to create a subscription: " + err.Error())
		return
	}

	sub.ReceiveSettings.Synchronous = true

	return
}

func (r *Reader) receiveFromPubsub(ctx context.Context) error {

	pubsubCallback := func(ctx context.Context, msg *pubsub.Message) {
		progressMsg, err := pubsubq.GetProgressMessage(msg)
		if err != nil {
			r.logger.Warn("can't unmarshall the message: " + err.Error())
			return
		}

		_, ok := r.progressPerReq[progressMsg.RequestID]
		if !ok {
			r.progressPerReq[progressMsg.RequestID] = make(chan int, maxProgress)
		}
		r.progressPerReq[progressMsg.RequestID] <- progressMsg.Progress

		if progressMsg.Progress >= maxProgress {
			close(r.progressPerReq[progressMsg.RequestID])
		}
	}

	// blocking call
	// receives messages on multiple goroutines
	r.logger.Debug("listening on progress queue...")
	return r.sub.Receive(ctx, pubsubCallback)
}
