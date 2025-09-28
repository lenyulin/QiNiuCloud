package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type ModelProviderResultProducer interface {
	AddEvent(evt AddEvent) error
}

const (
	MaxTCCMangerProduceRetry = 3
)

var (
	ErrProduceSubmitFailure = errors.New("failed to produce summit events")
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

var PartitionCount = 3

func (s *SaramaSyncProducer) AddEvent(evt AddEvent) error {
	val, err := json.Marshal(&evt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt.TimeStamp = time.Now().UnixMilli()
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicWatchEvent,
		Value: sarama.StringEncoder(val),
	})
	for err != nil {
		evt.Partition = partition
		evt.Topic = TopicWatchEvent
		evt.Offset = strconv.Itoa(int(offset))
		evt.Retry += 1
		if evt.Retry == MaxTCCMangerProduceRetry {
			//发送到死信队列
			return ErrProduceSubmitFailure
		}
	}
	return nil
}
