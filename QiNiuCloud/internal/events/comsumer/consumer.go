package consumer

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/internal/service"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/saramax"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

type InsertModelInfoConsumer struct {
	client sarama.Client
	log    logger.LoggerV1
	svc    service.ModelsService
}

func NewTCCManagerWatchEventConsumer(client sarama.Client, log logger.LoggerV1) *InsertModelInfoConsumer {
	return &InsertModelInfoConsumer{
		client: client,
		log:    log,
	}
}

func (i *InsertModelInfoConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("tcc_manager_group", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicInsertModelInfoToDBEvent},
			saramax.NewHandler[AddEvent](i.log, i.Consume))
		if er != nil {
			i.log.Error("Consume ERROR", logger.Error(er))
		}
	}()
	return nil
}

func (i *InsertModelInfoConsumer) Consume(msg *sarama.ConsumerMessage, event AddEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var model domain.ModelsInfo
	err := json.Unmarshal(msg.Value, &model)
	if err != nil {
		i.log.Error("Unmarshal ERROR", logger.Error(err))
		return err
	}
	return i.svc.AddModelToDB(ctx, model)
}
