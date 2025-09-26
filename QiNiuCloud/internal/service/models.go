package service

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/internal/repository"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/textshrink"
	"context"
)

type ModelsService interface {
	GenerateModel(ctx context.Context, text string) (string, error)
}

type service struct {
	l      logger.ZapLogger
	shrink textshrink.Shrink
	repo   repository.ModelsRepository
}

func (s *service) GenerateModel(ctx context.Context, text string) ([]domain.ModelsInfo, error) {
	keywords, err := s.shrink.Shrink(ctx, text)
	if err != nil {
		return nil, err
	}
	return s.repo.GetModelsByToken(ctx, keywords)
}
