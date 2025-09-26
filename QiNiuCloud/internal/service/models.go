package service

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/internal/repository"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/snowflake"
	"QiNiuCloud/QiNiuCloud/pkg/tcc"
	"QiNiuCloud/QiNiuCloud/pkg/tcc/event"
	"QiNiuCloud/QiNiuCloud/pkg/textshrink"
	"context"
	"errors"
	"strconv"
)

type ModelsService interface {
	GenerateModel(ctx context.Context, text string) (string, error)
}

type service struct {
	l                     logger.ZapLogger
	shrink                textshrink.Shrink
	repo                  repository.ModelsRepository
	snowflake             snowflake.Snowflake
	tccSaramaSyncProducer event.TCCMegProducer
}

const TccManagerProduceAddEvtTopic = "model_generate_tcc_evt"

func (s *service) GenerateModel(ctx context.Context, text string) ([]domain.ModelsInfo, string, error) {
	keywordstoken, err := s.shrink.Shrink(ctx, text)
	if err != nil {
		return nil, "", err
	}
	res, err := s.repo.GetModelsByToken(ctx, keywordstoken)
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
			evt := event.AddTCCEvent{
				TCCIdx: strconv.FormatInt(txid, 10),
				Topic:  TccManagerProduceAddEvtTopic,
				DATA: &tcc.ModelGenerateTransactionData{
					KeyWordsToken: keywordstoken,
					TransactionId: strconv.FormatInt(txid, 10),
				},
			}
			err = s.tccSaramaSyncProducer.TCCMangerProduceAddTCCEvent(evt)
			if err != nil {
				s.l.Error("Add Transaction Task failed",
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
