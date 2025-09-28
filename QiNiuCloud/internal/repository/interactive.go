package repository

import (
	"QiNiuCloud/QiNiuCloud/internal/repository/cache"
	"QiNiuCloud/QiNiuCloud/internal/repository/dao"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type InteractiveRepository interface {
	IncrLinkCnt(ctx context.Context, token string, hash string) error
	IncrDownloadCnt(ctx context.Context, token string, hash string) error
	IncrCloseAfterDownloadedCnt(ctx context.Context, token string, hash string) error
	//BatchIncrLinkCnt(ctx context.Context, bizs []string, bizIds []int64) error
	AddRecord(ctx context.Context, token string, hash string) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
	l     logger.LoggerV1
}

func NewCachedInteractiveRepository(dao dao.InteractiveDao, cache cache.InteractiveCache, l logger.LoggerV1) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

var (
	ErrInteractiveRecordNotFound = errors.New("interactive Record Not Found")
)

func (s *CachedInteractiveRepository) AddRecord(ctx context.Context, token string, hash string) error {
	//TODO implement me
	panic("implement me")
}
func (c *CachedInteractiveRepository) IncrCloseAfterDownloadedCnt(ctx context.Context, token string, hash string) error {
	err := c.cache.IncrCloseAfterDownloadedCntIfPresent(ctx, token, hash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//
			c.l.Error("interactive Record Not Found", logger.String("token", token), logger.String("hash", hash))
			return ErrInteractiveRecordNotFound
		}
		return err
	}
	return nil
}

func (c *CachedInteractiveRepository) IncrDownloadCnt(ctx context.Context, token string, hash string) error {
	err := c.cache.IncrDownloadCntIfPresent(ctx, token, hash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.l.Error("interactive Record Not Found", logger.String("token", token), logger.String("hash", hash))
			return ErrInteractiveRecordNotFound
		}
		return err
	}
	return nil
}

func (c *CachedInteractiveRepository) IncrLinkCnt(ctx context.Context, token string, hash string) error {
	err := c.cache.IncrLikeCntIfPresent(ctx, token, hash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//
			c.l.Error("interactive Record Not Found", logger.String("token", token), logger.String("hash", hash))
			return ErrInteractiveRecordNotFound
		}
		return err
	}
	return nil
}
