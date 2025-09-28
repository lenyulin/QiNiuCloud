package dao

import (
	"context"
)

type InteractiveDao interface {
	AddRecord(ctx context.Context, token string, model Models) error
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
}
