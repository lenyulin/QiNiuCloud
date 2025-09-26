package cache

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"context"
)

type Cache interface {
	FindByObjId(ctx context.Context, key string) ([]domain.ModelsInfo, error)
	Set(ctx context.Context, key string, value string) error
	GetModelByHash(ctx context.Context, hash string) (bool, error)
}
