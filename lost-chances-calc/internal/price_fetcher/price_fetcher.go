package pricefetcher

import (
	"context"
	"errors"
	"lost-chances-calc/internal/config"
	"lost-chances-calc/internal/domain"
	"sync"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

const (
	subscriptionName = "historicalPricesSub"
	topicName        = "historical-prices"
)

type PriceFetcher struct {
	logger *zap.Logger

	pricesSubscribed map[string]([]domain.HistoricalPrice) // key is the request ID
	sub              *pubsub.Subscription
	wg               sync.WaitGroup

	mutex sync.Mutex
}

func NewFetcher(logger *zap.Logger) PriceFetcher {
	return PriceFetcher{
		logger:           logger,
		wg:               sync.WaitGroup{},
		pricesSubscribed: make(map[string][]domain.HistoricalPrice),
		mutex:            sync.Mutex{},
	}
}

func (p *PriceFetcher) Init(ctx context.Context) (closer func(), err error) {
	cctx, ctxCancel := context.WithCancel(ctx)
	sub, subCloser, err := pubsubSubscribe(cctx)

	if err != nil {
		err = errors.New("subscribing to pubsub failed: " + err.Error())
		subCloser()
		ctxCancel()

		return
	}

	closer = func() {
		subCloser()
		ctxCancel()
		p.wg.Wait()
	}

	p.sub = sub

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := p.receiveFromPubsub(cctx); err != nil {
			p.logger.Error("error when receiving from pubsub: " + err.Error())
		}
		p.logger.Info("finished receiving from pubsub")

	}()

	return

}

func pubsubSubscribe(ctx context.Context) (sub *pubsub.Subscription, closerFunc func(), err error) {

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
		Topic:    client.Topic(topicName),
		Detached: false,
	}

	sub, err = client.CreateSubscription(ctx, subscriptionName, subCfg)
	if err != nil {
		err = errors.New("failed to create a subscription: " + err.Error())
		return
	}

	return
}

func (p *PriceFetcher) getHistoricalPrices(requestID string) []domain.HistoricalPrice {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.pricesSubscribed[requestID]
}

func (p *PriceFetcher) subscribeToHistoricalPrices(requestID string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.pricesSubscribed[requestID] = []domain.HistoricalPrice{}
}

func (p *PriceFetcher) unsubscribeFromHistoricalPrices(requestID string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.pricesSubscribed, requestID)
}

func (p *PriceFetcher) receiveFromPubsub(ctx context.Context) error {

	pubsubCallback := func(ctx context.Context, msg *pubsub.Message) {
		priceMsg, err := unmarshalHistoricalPriceMessage(msg.Data)
		if err != nil {
			p.logger.Warn(err.Error())
			return
		}

		p.logger.Debug("historical price read", zap.String("requestID", priceMsg.RequestID), zap.Any("price", priceMsg.Price))

		p.mutex.Lock()
		defer p.mutex.Unlock()

		list, ok := p.pricesSubscribed[priceMsg.RequestID]
		if !ok {
			p.logger.Debug("historical price not subscribed", zap.String("requestID", priceMsg.RequestID), zap.Any("price", priceMsg.Price))
			return
		}

		list = append(list, priceMsg.Price)
		p.pricesSubscribed[priceMsg.RequestID] = list
	}

	// blocking call
	// receives messages on multiple goroutines
	p.logger.Debug("listening on historical prices queue...")
	return p.sub.Receive(ctx, pubsubCallback)
}
