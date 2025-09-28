package cache

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"strings"
)

type ModelCache interface {
	FindByToken(ctx context.Context, key string) ([]domain.ModelsInfo, error)
	Set(ctx context.Context, key string, modelInfo domain.ModelsInfo) error
	GetModelByHash(ctx context.Context, token, hash string) (bool, error)
}

type redisCache struct {
	client redis.Cmdable
	l      logger.LoggerV1
}

func NewCache(client redis.Cmdable, l logger.LoggerV1) ModelCache {
	return &redisCache{
		client: client,
		l:      l,
	}
}
func (r *redisCache) FindByToken(ctx context.Context, key string) ([]domain.ModelsInfo, error) {
	res, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		r.l.Error("Record Not Found", logger.String("key", key), logger.Error(err))
		return nil, err
	}
	var result []domain.ModelsInfo
	for _, v := range res {
		var re domain.ModelsInfo
		err = json.Unmarshal([]byte(v), &re)
		result = append(result, re)
	}
	return result, nil
}

func (r *redisCache) Set(ctx context.Context, key string, modelInfo domain.ModelsInfo) error {
	value, err := json.Marshal(modelInfo)
	if err != nil {
		r.l.Error("Marshal Record Failed", logger.String("key", key), logger.Error(err))
		return err
	}
	_, err = r.client.SAdd(ctx, key, value).Result()
	if err != nil {
		r.l.Error("Add Record Failed", logger.String("key", key), logger.Error(err))
		return err
	}
	return nil
}

func (r *redisCache) GetModelByHash(ctx context.Context, token, hash string) (bool, error) {
	res, err := r.client.SMembers(ctx, token).Result()
	if err != nil {
		r.l.Error("Record Not Found", logger.String("key", token), logger.Error(err))
		return false, err
	}
	for _, v := range res {
		if strings.Contains(v, hash) {
			return true, nil
		}
	}
	return false, nil
}

func NewRedisCache(client redis.Cmdable, l logger.LoggerV1) ModelCache {
	return &redisCache{
		client: client,
		l:      l,
	}
}
