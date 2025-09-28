package ioc

import (
	"QiNiuCloud/QiNiuCloud/pkg/bloomfilterx"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
)

func InitBloomFilter(l logger.LoggerV1) bloomfilterx.BloomFilter {
	return bloomfilterx.NewbloomBuilder(l, 1000000, 0.01)
}
