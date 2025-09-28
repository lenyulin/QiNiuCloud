package ModelGnerationResultHelper

import (
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager/event"
	"QiNiuCloud/QiNiuCloud/pkg/ModelGnerationResultHelper/event/producer"
	"QiNiuCloud/QiNiuCloud/pkg/MultiOssManager"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"QiNiuCloud/QiNiuCloud/pkg/snowflake"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type ResultHelper interface {
	Process(ctx context.Context, jobId string, model ModelsInfo) error
}
type helper struct {
	l        logger.ZapLogger
	clients  *http.Client
	redis    *redis.Client
	sf       snowflake.Snowflake
	msgChan  chan *MultiOssManager.ActorMsg
	txTable  map[string]struct{}
	mu       sync.Mutex
	producer producer.ModelInfoInsertProducer
}

const (
	BASEDIR     = "./tmp_files/"
	MODELFORMAT = ".obj"
	THUMBNAIL   = ".jpg"
	MAXRETRYCNT = 3
)

func (h *helper) RunAsyncCheckUploadTx() {
	go func() {
		for {
			for k, _ := range h.txTable {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				values, err := h.redis.HMGet(ctx, k, "obj", "jpg", "token").Result()
				if err == nil {
					modelUrl := values[0].(string)
					thumbUrl := values[1].(string)
					token := values[2].(string)
					if modelUrl != "" && thumbUrl != "" && token != "" {
						er := h.writeToDB(modelUrl, thumbUrl, token)
						if er == nil {
							h.mu.Lock()
							delete(h.txTable, k)
							h.mu.Unlock()
						}
					}
				}
			}
		}
	}()
}
func (h *helper) writeToDB(modelUrl, thumbUrl, token string) error {
	d := &ModelsInfo{
		Token:     token,
		Url:       modelUrl,
		Thumbnail: thumbUrl,
	}
	err := h.producer.AddEvent(event.AddEvent(producer.AddEvent{
		DATA: d,
	}))
	if err != nil {
		return err
	}
	return nil
}
func (h *helper) Process(ctx context.Context, jobId string, model ModelsInfo) error {
	txId, _ := h.sf.NextID()
	tx := strconv.FormatInt(txId, 10)
	retry := 0
	err := errors.New("")
	var modelHash string
	for err != nil && retry <= MAXRETRYCNT {
		modelHash, err = h.downloadFile(BASEDIR, model.Url, MODELFORMAT)
		if err != nil {
			retry += 1
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return err
	}
	var thumbHash string
	err = errors.New("")
	retry = 0
	for err != nil && retry <= MAXRETRYCNT {
		thumbHash, err = h.downloadFile(BASEDIR, model.Thumbnail, THUMBNAIL)
		if err != nil {
			retry += 1
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		//进一步处理
		return err
	}
	err = h.redis.HSet(ctx, tx, "token", model.Token).Err()
	if err != nil {
		return err
	}
	h.mu.Lock()
	h.txTable[tx] = struct{}{}
	h.mu.Unlock()
	h.msgChan <- &MultiOssManager.ActorMsg{
		Op:       MultiOssManager.UploadFileToOSS,
		TxId:     tx,
		FileName: BASEDIR + modelHash + MODELFORMAT,
	}
	h.msgChan <- &MultiOssManager.ActorMsg{
		Op:       MultiOssManager.UploadFileToOSS,
		TxId:     tx,
		FileName: BASEDIR + thumbHash + THUMBNAIL,
	}
	return nil
}

func (h *helper) downloadFile(filepath string, url string, fileFmt string) (string, error) {
	resp, err := h.clients.Get(url)
	if err != nil {
		h.l.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.l.Error(err.Error())
		return "", err
	}
	hash := sha256.Sum256(body)
	sha256Sum := hex.EncodeToString(hash[:])
	filepath = filepath + sha256Sum + fileFmt
	err = os.WriteFile(filepath, body, 0644)
	if err != nil {
		h.l.Error(err.Error())
		return "", err
	}
	return sha256Sum, nil
}
