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
	Set(ctx context.Context, key string, model domain.ModelsInfo) error
	GetModelByHash(ctx context.Context, token, hash string) (bool, error)
}
type models struct {
	l                logger.ZapLogger
	modelCache       cache.ModelCache
	interactiveCache cache.InteractiveCache
	dao              dao.DAO
	filter           bloomfilterx.BloomFilter
}

func NewModelsRepository(l logger.ZapLogger, cache cache.ModelCache, dao dao.DAO, filter bloomfilterx.BloomFilter, interactiveCache cache.InteractiveCache) ModelsRepository {
	return &models{
		l:                l,
		modelCache:       cache,
		dao:              dao,
		filter:           filter,
		interactiveCache: interactiveCache,
	}
}
func (m *models) GetModelByHash(ctx context.Context, token, hash string) (bool, error) {
	found, _ := m.modelCache.GetModelByHash(ctx, token, hash)
	if found {
		return found, nil
	}
	return m.dao.GetModelByHash(ctx, token, hash)
}

func (m *models) GetModelsByToken(ctx context.Context, key string) ([]domain.ModelsInfo, error) {
	res, err := m.modelCache.FindByToken(ctx, key)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, cache.ErrRecordNotFound) || errors.Is(err, cache.ErrFailedConnectToCache) {
		if errors.Is(err, cache.ErrFailedConnectToCache) {
			m.l.Error("Error Failed Connect To ModelCache",
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
		go func() {
			for _, r := range res {
				err = m.modelCache.Set(ctx, key, r)
				if err != nil {
					m.l.Error("Set Model Cache Failed",
						logger.String("key", key),
						logger.String("hash", r.Hash),
						logger.Field{
							Key: "error",
							Val: err.Error()},
					)
				}
			}
		}()
		go func() {
			for _, r := range res {
				err = m.interactiveCache.Set(ctx, domain.Interactive{
					Token:                     r.Token,
					Hash:                      r.Hash,
					DownloadCount:             r.DownloadCount,
					LikeCount:                 r.LikeCount,
					CloseAfterDownloadedCount: r.CloseAfterDownloadedCount,
				})
				if err != nil {
					m.l.Error("Set Interactive Cache Failed",
						logger.String("key", key),
						logger.String("hash", r.Hash),
						logger.Field{
							Key: "error",
							Val: err.Error()},
					)
				}
			}
		}()
		return res, nil
	}
	m.l.Error("Unexpected Error",
		logger.Field{
			Key: "error",
			Val: err.Error()},
	)
	return nil, ErrInternalServerError
}

func (m *models) Set(ctx context.Context, key string, model domain.ModelsInfo) error {
	err := m.dao.Set(ctx, key, model)
	if err != nil {
		m.l.Error("Set Model Error", logger.String("key", key), logger.String("hash", model.Hash))
		return err
	}
	err = m.modelCache.Set(ctx, key, model)
	if err != nil {
		m.l.Error("Set Model Cache Error", logger.String("key", key), logger.String("hash", model.Hash))
		return err
	}
	err = m.interactiveCache.Set(ctx, domain.Interactive{
		Token:                     model.Token,
		Hash:                      model.Hash,
		DownloadCount:             model.DownloadCount,
		LikeCount:                 model.LikeCount,
		CloseAfterDownloadedCount: model.CloseAfterDownloadedCount,
	})
	if err != nil {
		m.l.Error("Set Interactive Cache Failed",
			logger.String("key", key),
			logger.String("hash", model.Hash))
		return err
	}
	return nil
}
