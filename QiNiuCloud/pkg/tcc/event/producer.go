package event

import (
	"QiNiuCloud/QiNiuCloud/pkg/tcc"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type TCCMegProducer interface {
	TCCMangerProduceAddTCCEvent(evt AddTCCEvent) error
}

const (
	MaxTCCMangerProduceRetry = 3
)

var (
	ErrProduceSubmitFailure = errors.New("failed to produce summit event")
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
	redis    *redis.Client
}

var PartitionCount = 3

func (s *SaramaSyncProducer) AddTCCEvent(evt AddTCCEvent) error {
	evt.Status = TransactionStatus(tcc.StatusTrying)
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
