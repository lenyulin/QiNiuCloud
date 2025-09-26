package repository

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/internal/repository/cache"
	"QiNiuCloud/QiNiuCloud/internal/repository/dao"
	"QiNiuCloud/QiNiuCloud/pkg/bloomfilterx"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"errors"
)

type ModelsRepository interface {
	GetModelsByToken(ctx context.Context, key string) ([]domain.ModelsInfo, error)
	Set(ctx context.Context, key string, value string) error
	GetModelByHash(ctx context.Context, token, hash string) (bool, error)
}
type models struct {
	l      logger.ZapLogger
	cache  cache.Cache
	dao    dao.DAO
	filter bloomfilterx.BloomFilter
}

func NewModelsRepository(l logger.ZapLogger, cache cache.Cache, dao dao.DAO, filter bloomfilterx.BloomFilter) ModelsRepository {
	return &models{
		l:      l,
		cache:  cache,
		dao:    dao,
		filter: filter,
	}
}
func (m *models) GetModelByHash(ctx context.Context, token, hash string) (bool, error) {
	found, _ := m.cache.GetModelByHash(ctx, token, hash)
	if found {
		return found, nil
	}
	return m.dao.GetModelByHash(ctx, token, hash)
}

func (m *models) GetModelsByToken(ctx context.Context, key string) ([]domain.ModelsInfo, error) {
	res, err := m.cache.FindByToken(ctx, key)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, cache.ErrRecordNotFound) || errors.Is(err, cache.ErrFailedConnectToCache) {
		if errors.Is(err, cache.ErrFailedConnectToCache) {
			m.l.Error("Error Failed Connect To Cache",
				logger.Field{
					Key: "error",
					Val: err.Error()},
			)
			return nil, err
		}
		//查询过滤器
		find, err := m.filter.Get(ctx, key)
		if err != nil {
			m.l.Error("Bloom Filter Error",
				logger.Field{
					Key: "error",
					Val: err.Error()},
			)
			return nil, err
		}
		if !find {
			m.l.Debug("Bloom Filter Not Found Key",
				logger.Field{
					Key: "Debug",
					Val: ErrBloomFilterNotFoundRecord},
			)
			return nil, ErrBloomFilterNotFoundRecord
		}
		res, err = m.dao.FindByObjId(ctx, key)
		if err != nil {
			m.l.Error("Record Not Found In DataBase",
				logger.Field{
					Key: "error",
					Val: err.Error()},
			)
			return nil, ErrResourceNotFound
		}
		return res, nil
	}
	m.l.Error("Unexpected Error",
		logger.Field{
			Key: "error",
			Val: err.Error()},
	)
	return nil, ErrInternalServerError
}

func (m *models) Set(ctx context.Context, key string, value string) error {
	//TODO implement me
	panic("implement me")
}
