package service

import (
	"QiNiuCloud/QiNiuCloud/internal/repository"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/textshrink"
	"context"
)

type ModelsService interface {
	Generate(ctx context.Context, text string) (string, error)
}

type service struct {
	l      logger.ZapLogger
	shrink textshrink.Shrink
	repo   repository.ModelsRepository
}

func (s *service) Generate(ctx context.Context, text string) (string, error) {
	keywords, err := s.shrink.Shrink(ctx, text)
	if err != nil {
		return "", err
	}
	return s.repo.Get(ctx, keywords)
}
