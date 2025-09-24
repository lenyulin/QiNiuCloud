package repository

import (
	"QiNiuCloud/QiNiuCloud/internal/repository/cache"
	"QiNiuCloud/QiNiuCloud/internal/repository/dao"
	"QiNiuCloud/QiNiuCloud/pkg/bloomfilterx"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"errors"
)

type ModelsRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}
type models struct {
	l      logger.ZapLogger
	cache  cache.Cache
	dao    dao.DAO
	filter bloomfilterx.BloomFilter
}

func (m *models) Get(ctx context.Context, key string) (string, error) {
	res, err := m.cache.Get(ctx, key)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, cache.ErrRecordNotFound) {
		//查询过滤器
		find, err := m.filter.Get(ctx, key)
		if err != nil {
			m.l.Error("Bloom Filter Error",
				logger.Field{
					Key: "error",
					Val: err.Error()},
			)
			return "", err
		}
		if !find {
			m.l.Debug("Bloom Filter Not Found Key",
				logger.Field{
					Key: "Debug",
					Val: ErrBloomFilterNotFoundRecord},
			)
			return "", ErrBloomFilterNotFoundRecord
		}
		//查询MySQL
		//TODO implement me
		panic("implement me")
	}
	return "", err
}

func (m *models) Set(ctx context.Context, key string, value string) error {
	//TODO implement me
	panic("implement me")
}
