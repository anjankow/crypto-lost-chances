package progressupdates

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lost-chances-calc/internal/config"
	"sync"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

const (
	topicName = "progress-update"

	maxProgress = 100
	minProgress = 0
)

type Writer struct {
	logger *zap.Logger
	topic  *pubsub.Topic
	wg     *sync.WaitGroup
}

func NewWriter(logger *zap.Logger) Writer {
	return Writer{
		logger: logger,
		wg:     &sync.WaitGroup{},
	}
}

func (w *Writer) Init(ctx context.Context) (closer func(), err error) {
	cctx, ctxCancel := context.WithCancel(ctx)

	var client *pubsub.Client

	closer = func() {
		ctxCancel()
		client.Close()
	}

	client, err = pubsub.NewClient(cctx, config.GetProjectID())
	if err != nil {
		err = errors.New("failed to create a pubsub client: " + err.Error())
		return
	}

	topic := client.Topic(topicName)
	if topic == nil {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			err = errors.New("failed to create a new topic: " + err.Error())
			return
		}
		topic.EnableMessageOrdering = true

	}

	w.topic = topic

	return

}

func (w Writer) PublishProgress(ctx context.Context, requestID string, progress int) error {
	if progress > maxProgress || progress < minProgress {
		return fmt.Errorf("invalid progress value, shall be in range (0,100): ", progress)
	}

	msgBody := progressMessage{
		RequestID: requestID,
		Progress:  progress,
	}
	bytes, err := json.Marshal(msgBody)
	if err != nil {
		return errors.New("error when marshalling the message: " + err.Error())
	}

	result := w.topic.Publish(ctx, &pubsub.Message{
		Data: bytes,
	})

	w.wg.Add(1)
	go func(result *pubsub.PublishResult) {
		defer w.wg.Done()

		_, err := result.Get(ctx)
		if err != nil {
			w.logger.Error("error when publishing the progress message: " + err.Error())
		}
	}(result)

	return nil
}
