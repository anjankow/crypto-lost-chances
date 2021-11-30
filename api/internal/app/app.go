package app

import (
	"api/internal/pubsubq"
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	sub    *pubsub.Subscription
}

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

type Results struct {
	Cryptocurrency string
	Income         float32
}

func NewApp(l *zap.Logger) (app App, closer func(), err error) {
	sub, closer, err := pubsubq.SubscribeToProgressUpdates(context.Background())
	if err != nil {
		l.Error(err.Error())
		closer()
		closer = func() {}

		sub = nil
	}

	app = App{
		Logger: l,
		sub:    sub,
	}
	return
}

func (a App) ProcessCalculateReq(ctx context.Context, input UserInput) (Results, error) {
	results := Results{Cryptocurrency: "ADA", Income: float32(input.Amount * 2)}

	// calls the main app

	return results, nil
}

// ListenProgress listens on the queue for the request progress
func (a App) ListenProgress(ctx context.Context, requestID string, callback func(progress int)) error {

	if a.sub == nil {
		return errors.New("not subscribed to the queue")
	}

	pubsubCallback := func(ctx context.Context, msg *pubsub.Message) {
		progressMsg, err := pubsubq.GetProgressMessage(msg)
		if err != nil {
			a.Logger.Warn("can't unmarshall the message: " + err.Error())
			return
		}

		if progressMsg.RequestID != requestID {
			a.Logger.Debug("received pubsub message for another request", zap.String("expected", strings.Split(requestID, "-")[0]), zap.String("received", strings.Split(progressMsg.RequestID, "-")[0]))
			return
		}

		callback(progressMsg.Progress)
	}

	defer func() {
		a.Logger.Debug("finished listening...", zap.String("requestID", requestID))
	}()

	// blocking call
	// receives messages on multiple goroutines
	a.Logger.Debug("listening on progress queue...", zap.String("requestID", requestID))
	return a.sub.Receive(ctx, pubsubCallback)
}
