package service

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/internal/repository"
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/snowflake"
	"QiNiuCloud/QiNiuCloud/pkg/textshrink"
	"context"
	"errors"
	"strconv"
)

type ModelsService interface {
	GenerateModel(ctx context.Context, text string) (string, error)
	AddModelToDB(ctx context.Context, model domain.ModelsInfo) error
}

type service struct {
	l                logger.ZapLogger
	shrink           textshrink.Shrink
	modelRepo        repository.ModelsRepository
	snowflake        snowflake.Snowflake
	generatorManager AsyncModelGenerationTaskManager.SyncModelGenerationTaskManager
}

const TccManagerProduceAddEvtTopic = "model_generate_tcc_evt"

var (
	ErrAddModelToDB    = errors.New("add model to DB error")
	ErrAddModelToCache = errors.New("add model to Cache error")
)

func (s *service) AddModelToDB(ctx context.Context, model domain.ModelsInfo) error {
	err := s.modelRepo.Set(ctx, model.Token, model)
	if err != nil {
		s.l.Error("Add Model To DB failed",
			logger.Field{
				Key: "error",
				Val: err.Error()},
		)
		return ErrAddModelToDB
	}
	return nil
}

func (s *service) GenerateModel(ctx context.Context, text string) ([]domain.ModelsInfo, string, error) {
	keywordstoken, err := s.shrink.Shrink(ctx, text)
	if err != nil {
		return nil, "", err
	}
	res, err := s.modelRepo.GetModelsByToken(ctx, keywordstoken)
	if err != nil {
		if errors.Is(err, repository.ErrResourceNotFound) {
			txid, err := s.snowflake.NextID()
			if err != nil {
				s.l.Error("Generate TX id failed",
					logger.Field{
						Key: "error",
						Val: err.Error()},
				)
				return nil, "", ErrResourceNotFound
			}
			err = s.generatorManager.AddTask(ctx, strconv.FormatInt(txid, 10), keywordstoken)
			if err != nil {
				s.l.Error("Add Model Generation Task failed",
					logger.Field{
						Key: "error",
						Val: err.Error()},
				)
				return nil, "", ErrResourceNotFound
			}
			return nil, strconv.FormatInt(txid, 10), ErrResourceNotFound
		}
		return nil, "", err
	}
	return res, "", nil
}
