package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/rwpp/RzWeLook/interactive/repository"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/rwpp/RzWeLook/pkg/saramax"
	"time"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func (k *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("read_event",
		k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_event"},
			saramax.NewBatchHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("消费消息失败", logger.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventBatchConsumer) Consume(msg []*sarama.ConsumerMessage, t []ReadEvent) error {
	ids := make([]int64, 0, len(t))
	bizs := make([]string, 0, len(t))
	for _, evt := range t {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err := k.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		k.l.Error("批量消费消息失败",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil
}
