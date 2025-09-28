package producer

import (
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager/event"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type ModelInfoInsertProducer interface {
	AddEvent(evt event.AddEvent) error
}

const (
	MaxProduceRetry = 3
)

var (
	ErrProduceSubmitFailure = errors.New("failed to produce summit events")
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

var PartitionCount = 3

func (s *SaramaSyncProducer) AddEvent(evt event.AddEvent) error {
	val, err := json.Marshal(&evt.DATA)
	if err != nil {
		fmt.Println(err)
		return err
	}
	evt.TimeStamp = time.Now().UnixMilli()
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicInsertModelInfoToDBEvent,
		Value: sarama.StringEncoder(val),
	})
	for err != nil {
		evt.Partition = partition
		evt.Topic = TopicInsertModelInfoToDBEvent
		evt.Offset = strconv.Itoa(int(offset))
		evt.Retry += 1
		if evt.Retry == MaxProduceRetry {
			//发送到死信队列
			return ErrProduceSubmitFailure
		}
	}
	return nil
}
