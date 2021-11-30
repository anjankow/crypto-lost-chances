package app

import (
	"api/internal/config"
	"api/internal/pubsubq"
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

const (
	progressTopic = "progress-update"
)

type App struct {
	Logger   *zap.Logger
	psClient *pubsub.Client
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
	projectID := config.GetProjectID()

	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		err = errors.New("failed to create a pubsub client: " + err.Error())
		return
	}

	return App{
		Logger:   l,
		psClient: client,
	}, func() { client.Close() }, nil
}

func (a App) ProcessCalculateReq(ctx context.Context, input UserInput) (Results, error) {
	results := Results{Cryptocurrency: "ADA", Income: float32(input.Amount * 2)}

	// calls the main app

	return results, nil
}

// ListenProgress listens on the queue for the request progress
func (a App) ListenProgress(ctx context.Context, requestID string, callback func(progress int)) error {

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

	sub, err := pubsubq.Subscribe(ctx, a.psClient)
	if err != nil {
		return err
	}

	defer func() {
		a.Logger.Debug("finished listening...", zap.String("requestID", requestID))
	}()

	// blocking call
	// receives messages on multiple goroutines
	a.Logger.Debug("listening on progress queue...", zap.String("requestID", requestID))
	return sub.Receive(ctx, pubsubCallback)
}
