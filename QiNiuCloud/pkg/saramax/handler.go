package saramax

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	log logger.LoggerV1
	fn  func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](log logger.LoggerV1, fn func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {
	return &Handler[T]{
		log: log,
		fn:  fn,
	}
}
func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.log.Error("Unmarshal Failed",
				logger.String("Topic", msg.Topic),
				logger.Int32("Partition", msg.Partition),
				logger.String("Key", string(msg.Key)),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		err = h.fn(msg, t)
		if err != nil {
			h.log.Error("Consume Failed",
				logger.String("Topic", msg.Topic),
				logger.Int32("Partition", msg.Partition),
				logger.String("Key", string(msg.Key)),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
