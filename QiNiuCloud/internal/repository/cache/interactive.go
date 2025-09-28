package cache

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
	//go:embed lua/set_interactive.lua
	setInteractive string
)

const (
	fieldLikeCntCnt              = "like_cnt"
	fieldDownloadCntCnt          = "download_cnt"
	fieldCloseAfterDownloadedCnt = "close_after_downloaded_cnt"
)

type InteractiveCache interface {
	IncrLikeCntIfPresent(ctx context.Context, token string, hash string) error
	IncrDownloadCntIfPresent(ctx context.Context, token string, hash string) error
	IncrCloseAfterDownloadedCntIfPresent(ctx context.Context, token string, hash string) error
	Set(ctx context.Context, interactive domain.Interactive) error
}
type InteractiveRedisCache struct {
	l      logger.LoggerV1
	client redis.Cmdable
}

func NewInteractiveRedisCache(l logger.LoggerV1, client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		l:      l,
		client: client,
	}
}
func (i *InteractiveRedisCache) IncrCloseAfterDownloadedCntIfPresent(ctx context.Context, token string, hash string) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.generateKey(token, hash)}, fieldCloseAfterDownloadedCnt, 1).Err()
}
func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context, token string, hash string) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.generateKey(token, hash)}, fieldLikeCntCnt, 1).Err()
}
func (i *InteractiveRedisCache) IncrDownloadCntIfPresent(ctx context.Context, token string, hash string) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.generateKey(token, hash)}, fieldDownloadCntCnt, 1).Err()
}
func (i *InteractiveRedisCache) Set(ctx context.Context, interactive domain.Interactive) error {
	args := []interface{}{
		fieldLikeCntCnt, fieldDownloadCntCnt, fieldCloseAfterDownloadedCnt,
		interactive.LikeCount, interactive.DownloadCount, interactive.CloseAfterDownloadedCount}
	return i.client.Eval(ctx, setInteractive,
		[]string{i.generateKey(interactive.Token, interactive.Hash)}, args).Err()
}
func (i *InteractiveRedisCache) generateKey(token string, hsah string) string {
	return fmt.Sprintf("interative:model:%s:%d", token, hsah)
}
