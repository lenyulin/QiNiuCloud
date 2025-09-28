package consumer

import (
	"QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/saramax"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

type ManagerWatchEventConsumer struct {
	client sarama.Client
	log    logger.LoggerV1
	helper ModelGnerationResultHelper.ResultHelper
}

func NewTCCManagerWatchEventConsumer(client sarama.Client, log logger.LoggerV1) *ManagerWatchEventConsumer {
	return &ManagerWatchEventConsumer{
		client: client,
		log:    log,
	}
}

func (i *ManagerWatchEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("tcc_manager_group", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicWatchEvent},
			saramax.NewHandler[AddEvent](i.log, i.Consume))
		if er != nil {
			i.log.Error("TCC Consume ERROR", logger.Error(er))
		}
	}()
	return nil
}

func (i *ManagerWatchEventConsumer) Consume(msg *sarama.ConsumerMessage, event AddEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var model ModelGnerationResultHelper.ModelGenerationTaskResult
	err := json.Unmarshal(msg.Value, &model)
	if err != nil {
		i.log.Error("Unmarshal ERROR", logger.Error(err))
		return err
	}
	return i.helper.Process(ctx, model.JobId, ModelGnerationResultHelper.ModelsInfo{
		Token:     model.Token,
		Url:       model.Url,
		Thumbnail: model.Thumb,
	})
}
