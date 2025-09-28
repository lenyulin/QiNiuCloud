package ioc

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/textshrink"
	"net/http"
)

func InitShrink(clientPool *http.Client, l logger.LoggerV1) textshrink.Shrink {
	return textshrink.NewShrink(clientPool, l)
}
