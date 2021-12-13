package progressupdates

import (
	"api/internal/config"
	"context"
	"errors"
	"sync"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

const (
	subscriptionName = "progressSub"
	topicName        = "progress-update"

	channelLen = 100
)

type Reader struct {
	logger *zap.Logger
	sub    *pubsub.Subscription
	// this approach could be improved - channels should be deleted when not needed anymore
	// no need to keep them in memory after request
	progressPerReq map[string]chan int
	wg             sync.WaitGroup
}

func NewReader(logger *zap.Logger) Reader {
	return Reader{
		logger:         logger,
		progressPerReq: map[string]chan int{},
		wg:             sync.WaitGroup{},
	}
}

func (r *Reader) Subscribe(requestID string) chan int {
	_, ok := r.progressPerReq[requestID]
	if !ok {
		r.progressPerReq[requestID] = make(chan int, channelLen)
	}

	return r.progressPerReq[requestID]

}

func (r *Reader) Start(ctx context.Context) (closer func(), err error) {

	cctx, ctxCancel := context.WithCancel(ctx)
	sub, subCloser, err := subscribePubsub(cctx)
	closer = func() {
		subCloser()
		ctxCancel()
		r.wg.Wait()
	}

	if err != nil {
		err = errors.New("subscribing to pubsub failed: " + err.Error())
		return
	}

	r.sub = sub

	r.wg.Add(1)
	go func() {
		if err := r.receiveFromPubsub(cctx); err != nil {
			r.logger.Error("error when receiving from pubsub: " + err.Error())
		}
		r.logger.Info("finished receiving from pubsub")
		r.wg.Done()
	}()

	return

}

func subscribePubsub(ctx context.Context) (sub *pubsub.Subscription, closerFunc func(), err error) {

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
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	return
}

func (r *Reader) receiveFromPubsub(ctx context.Context) error {

	pubsubCallback := func(ctx context.Context, msg *pubsub.Message) {
		progressMsg, err := unmarshalProgressMessage(msg)
		if err != nil {
			r.logger.Warn("can't unmarshall the message: " + err.Error())
			return
		}

		r.logger.Debug("progress read", zap.String("requestID", progressMsg.RequestID), zap.Int("progress", progressMsg.Progress))

		_, ok := r.progressPerReq[progressMsg.RequestID]

		if !ok {

			// the progress messages might start coming before subscription
			r.progressPerReq[progressMsg.RequestID] = make(chan int, channelLen)
		}
		r.progressPerReq[progressMsg.RequestID] <- progressMsg.Progress
	}

	// blocking call
	// receives messages on multiple goroutines
	r.logger.Debug("listening on progress queue...")
	return r.sub.Receive(ctx, pubsubCallback)
}

func (r *Reader) Unsubscribe(requestID string) {
	r.logger.Debug("closing the channel", zap.String("requestID", requestID))
	close(r.progressPerReq[requestID])
	delete(r.progressPerReq, requestID)
}
