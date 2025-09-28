package ioc

import (
	consumer "QiNiuCloud/QiNiuCloud/internal/events/comsumer"
	"QiNiuCloud/QiNiuCloud/internal/service"
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager/event"
	"QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper"
	consumer2 "QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper/event/consumer"
	producerx "QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper/event/producer"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitConsumers(client sarama.Client, l logger.LoggerV1, svc service.ModelsService, helper ModelGnerationResultHelper.ResultHelper) []consumer.Consumer {
	return []consumer.Consumer{
		consumer.NewInsertModelInfoToDBEventConsumer(client, l, svc),
		consumer2.NewTCCManagerWatchEventConsumer(client, l, helper),
	}
}

// ModelProviderResultProducer
func InitModelProviderResultProducer(p sarama.SyncProducer) event.ModelProviderResultProducer {
	return event.NewModelProviderResultProducer(p)
}
func InitModelInfoInsertProducer(producer sarama.SyncProducer) producerx.ModelInfoInsertProducer {
	return producerx.NewModelInfoInsertProducer(producer)
}
func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}
func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string
	}
	var cfg Config
	cfg.Addr = []string{"127.0.0.1:9094"}
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}
