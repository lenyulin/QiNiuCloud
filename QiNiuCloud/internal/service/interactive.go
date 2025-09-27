package service

import (
	"QiNiuCloud/QiNiuCloud/internal/repository"
	"context"
)

type InteractiveService interface {
	IncrLinkCnt(ctx context.Context, token string, hash string) error
}
type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (s *interactiveService) IncrLinkCnt(ctx context.Context, token string, hash string) error {
	return s.repo.IncrLinkCnt(ctx, token, hash)
}
