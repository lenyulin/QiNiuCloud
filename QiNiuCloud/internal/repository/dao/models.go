package dao

import (
	"QiNiuCloud/QiNiuCloud/internal/domain"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DAO interface {
	FindByObjId(ctx context.Context, key string) ([]domain.ModelsInfo, error)
	Set(ctx context.Context, key string, model domain.ModelsInfo) error
	GetModelByHash(ctx context.Context, token, hash string) (bool, error)
}
type MongoDBDAO struct {
	client *mongo.Client
	l      logger.ZapLogger
}

const (
	DATABASE   = "models"
	COLLECTION = "collection"
)

func (dao *MongoDBDAO) GetModelByHash(ctx context.Context, token, hash string) (bool, error) {
	coll := dao.client.Database(DATABASE).Collection(COLLECTION)
	filter := bson.D{
		bson.E{"token", token},
		bson.E{"hash", hash},
	}
	var result Models
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, ErrRecordNotFound
		}
		dao.l.Error("UNKNOWN ERROR WHEN SEARCHING TOKEN: " + token)
		return false, ErrUnknown
	}
	return true, nil
}
func (dao *MongoDBDAO) FindByObjId(ctx context.Context, token string) ([]domain.ModelsInfo, error) {
	coll := dao.client.Database(DATABASE).Collection(COLLECTION)
	filter := bson.D{{"token", token}}
	var result Models
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		dao.l.Error("UNKNOWN ERROR WHEN SEARCHING TOKEN: " + token)
		return nil, ErrUnknown
	}
	if len(result.Context) == 0 {
		dao.l.Error("BAD MODEL RECORD WHEN SEARCHING TOKEN: " + token)
		return nil, ErrBadModelRecord
	}
	return dao.toDomain(result), nil
}
func (dao *MongoDBDAO) toDomain(models Models) []domain.ModelsInfo {
	var results []domain.ModelsInfo
	for _, m := range models.Context {
		results = append(results, domain.ModelsInfo{
			Thumbnail:     m.Thumbnail,
			Url:           m.Url,
			DownloadCount: m.DownloadCount,
			LikeCount:     m.LikeCount,
		})
	}
	return results
}
func (dao *MongoDBDAO) Set(ctx context.Context, key string, model domain.ModelsInfo) error {
	//TODO implement me
	panic("implement me")
}
