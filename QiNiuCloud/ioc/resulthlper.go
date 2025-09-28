package ioc

import (
	"QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper"
	"QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper/event/producer"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"github.com/redis/go-redis/v9"
	"net/http"
)

//NewResultHelper(l logger.LoggerV1, clients *http.Client, redis redis.Cmdable, producer producer.ModelInfoInsertProducer) ResultHelper

func InitResultHelper(l logger.LoggerV1, clients *http.Client, redis redis.Cmdable, producer producer.ModelInfoInsertProducer) ModelGnerationResultHelper.ResultHelper {
	return ModelGnerationResultHelper.NewResultHelper(l, clients, redis, producer)
}
