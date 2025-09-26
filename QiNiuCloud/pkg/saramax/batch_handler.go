package saramax

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

type BatchHandler[T any] struct {
	fn  func(msg []*sarama.ConsumerMessage, ts []T) error
	log logger.LoggerV1
}

func NewBatchHandler[T any](log logger.LoggerV1, fn func(msg []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{
		fn:  fn,
		log: log,
	}
}
func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchsize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchsize)
		ts := make([]T, 0, batchsize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var done = false
		for i := 0; i < batchsize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.log.Error("Unmarshal Failed @wedy/pkg/tdd/saramax/batch_handler.go line:50",
						logger.String("Topic", msg.Topic),
						logger.Int32("Partition", msg.Partition),
						logger.String("Key", string(msg.Key)),
						logger.Int64("offset", msg.Offset),
						logger.Error(err))
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		err := b.fn(batch, ts)
		if err != nil {
			b.log.Error("Consume Failed  @wedy/pkg/tdd/saramax/batch_handler.go line:64",
				logger.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
