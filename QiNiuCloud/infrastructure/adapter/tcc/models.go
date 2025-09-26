package tcc

import (
	"QiNiuCloud/QiNiuCloud/internal/service"
	"QiNiuCloud/QiNiuCloud/pkg/tcc"
	"context"
)

type ModelsAdapter struct {
	modelsService service.ModelsService // 依赖库存服务
}

func NewStockAdapter(modelsService service.ModelsService) tcc.TccAction {
	return &ModelsAdapter{
		modelsService: modelsService,
	}
}

func (s *ModelsAdapter) Try(ctx context.Context, tccID string, bizData interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s *ModelsAdapter) Confirm(ctx context.Context, tccID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *ModelsAdapter) Cancel(ctx context.Context, tccID string) error {
	//TODO implement me
	panic("implement me")
}
