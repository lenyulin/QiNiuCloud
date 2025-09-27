package repository

import (
	"QiNiuCloud/QiNiuCloud/internal/repository/cache"
	"QiNiuCloud/QiNiuCloud/internal/repository/dao"
	"context"
)

type InteractiveRepository interface {
	IncrLinkCnt(ctx context.Context, token string, hash string) error
	//BatchIncrLinkCnt(ctx context.Context, bizs []string, bizIds []int64) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
}

func NewCachedInteractiveRepository(dao dao.InteractiveDao, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
	}
}

func (c *CachedInteractiveRepository) IncrLinkCnt(ctx context.Context, token string, hash string) error {
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}
