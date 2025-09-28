package MultiOssManager

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	oss "QiNiuCloud/QiNiuCloud/pkg/ossx"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"path/filepath"
	"sync/atomic"
	"time"
)

type MultiOssManager interface {
	Upload(tx string, filename string)
	Delete()
	Find()
}
type OSSObject struct {
	Client        oss.OSSHandler // OSS客户端
	ObjectId      int
	IsAvailable   atomic.Bool
	IsFailure     atomic.Bool
	FailureCount  int       // 失效次数
	RetryCount    int       // 重试次数
	LastFailureAt time.Time // 最近失效时间
	LastModified  time.Time // 最近修改时间
}
type OSSPool struct {
	mainPool       []*OSSObject
	mainPoolSize   int
	backupPool     []*OSSObject
	backupPoolSize int
	failureOssPool map[int]*OSSObject
}

const (
	UploadFileToOSS = iota
	DeleteFileFromOSS
	FindFileFromOSS
	ChangeStatusToAvailable
	ChangeOssStatusToSuccess
	ChangeOssStatusToFail
)

type ActorMsg struct {
	ossObject *OSSObject
	Op        int
	FileName  string
	TxId      string
	status    Status
}
type Status struct {
	IsAvailable bool
	IsFailure   bool
}
type ossManager struct {
	l             logger.ZapLogger
	redis         *redis.Client
	nAvailableOss int64
	OSSPool       *OSSPool
	msgChan       chan *ActorMsg
}

func (o *ossManager) GetAvailableOssCount() int64 {
	return o.nAvailableOss
}
func (o *ossManager) RunActor() {
	for msg := range o.msgChan {
		switch msg.Op {
		case UploadFileToOSS:
			o.Upload(msg.TxId, msg.FileName)
		case DeleteFileFromOSS:
			o.Delete()
		case FindFileFromOSS:
			o.Find()
		case ChangeOssStatusToFail:
			o.markAsFailed(msg.ossObject)
		default:
			panic("unhandled default case")
		}
	}
}

func (o *ossManager) StartFailureOssManager() {
	for {
		for _, ossObject := range o.OSSPool.mainPool {
			if !ossObject.IsAvailable.Load() && ossObject.IsFailure.Load() {
				if o.testClientUpload(ossObject) == nil {
					ossObject.IsFailure.Store(false)
					atomic.AddInt64(&o.nAvailableOss, 1)
					ossObject.IsAvailable.Store(true)
				}
			}
		}
		for _, ossObject := range o.OSSPool.backupPool {
			if !ossObject.IsAvailable.Load() && ossObject.IsFailure.Load() {
				if o.testClientUpload(ossObject) == nil {
					ossObject.IsFailure.Store(false)
					atomic.AddInt64(&o.nAvailableOss, 1)
					ossObject.IsAvailable.Store(true)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

const TestFilePath = "./testFile.txt"

func (o *ossManager) testClientUpload(ossObj *OSSObject) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _, er := ossObj.Client.Upload(ctx, TestFilePath)
	if er != nil {
		o.l.Error(er.Error())
		return er
	}
	return nil
}

func (o *ossManager) markAsFailed(object *OSSObject) {
	object.FailureCount += 1
	object.LastFailureAt = time.Now()
	object.LastModified = time.Now()
	object.IsFailure.Store(true)
}

//	func (o *ossManager) markAsAvailable(object *OSSObject) {
//		object.IsAvailable.Swap(true)
//		object.LastModified = time.Now()
//	}
func (o *ossManager) getAvailableOss() (*OSSObject, error) {
	for _, obj := range o.OSSPool.mainPool {
		if obj != nil && obj.IsAvailable.Load() {
			atomic.AddInt64(&o.nAvailableOss, -1)
			return obj, nil
		}
	}
	for _, obj := range o.OSSPool.backupPool {
		if obj != nil && obj.IsAvailable.Load() {
			atomic.AddInt64(&o.nAvailableOss, -1)
			return obj, nil
		}
	}
	return nil, errors.New("ossManager getAvailableOssClient failed")
}
func NewOSSManager(l logger.ZapLogger, mainOss []oss.OSSHandler, backupOss []oss.OSSHandler) MultiOssManager {
	ch := make(chan *ActorMsg, 1000)
	ossPool := &OSSPool{
		mainPool:   make([]*OSSObject, 0),
		backupPool: make([]*OSSObject, 0),
	}
	for idx, o := range mainOss {
		if o != nil {
			ossObj := &OSSObject{
				Client:        o,
				ObjectId:      idx,
				FailureCount:  0,
				RetryCount:    0,
				LastFailureAt: time.Time{},
				LastModified:  time.Now(),
			}
			ossObj.IsAvailable.Store(true)
			ossObj.IsFailure.Store(false)
			ossPool.mainPool = append(ossPool.mainPool, ossObj)
			ossPool.mainPoolSize += 1
		}
	}
	for idx, o := range backupOss {
		if o != nil {
			ossObj := &OSSObject{
				Client:        o,
				ObjectId:      idx,
				FailureCount:  0,
				RetryCount:    0,
				LastFailureAt: time.Time{},
				LastModified:  time.Now(),
			}
			ossObj.IsAvailable.Store(true)
			ossObj.IsFailure.Store(false)
			ossPool.backupPool = append(ossPool.backupPool, ossObj)
			ossPool.backupPoolSize += 1
		}
	}
	return &ossManager{
		l:             l,
		nAvailableOss: int64(ossPool.backupPoolSize + ossPool.mainPoolSize),
		OSSPool:       ossPool,
		msgChan:       ch,
	}
}

const MaxRetryCount = 3

func (o *ossManager) Upload(tx string, filename string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	availableOss, err := o.getAvailableOss()
	if err != nil {
		//处理？
	}
	go func(ctx context.Context, oss *OSSObject) {
		retry := 0
		url := ""
		_, _, er := oss.Client.Upload(ctx, filename)
		for er != nil && retry <= MaxRetryCount {
			//进一步处理?
			retry += 1
			time.Sleep(time.Millisecond * 30)
			url, _, er = oss.Client.Upload(ctx, filename)
		}
		if er != nil {
			if errors.Is(err, ErrOssClientFailed) {
				o.msgChan <- &ActorMsg{
					Op:        ChangeOssStatusToFail,
					ossObject: oss,
					FileName:  filename,
					status: Status{
						IsAvailable: false,
						IsFailure:   true,
					},
				}
			}
			o.l.Error(err.Error())
			return
		}
		o.uploadDone(tx, filename, url)
		oss.IsAvailable.Store(true)
		atomic.AddInt64(&o.nAvailableOss, 1)
	}(ctx, availableOss)
}
func (o *ossManager) uploadDone(tx string, filename string, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	ext := filepath.Ext(filename)
	_, er := o.redis.HSet(ctx, tx, ext[1:], url).Result()
	if er != nil {
		//进一步处理
		o.l.Error(er.Error())
	}
}
func (o *ossManager) Delete() {
	//TODO implement me
	panic("implement me")
}

func (o *ossManager) Find() {
	//TODO implement me
	panic("implement me")
}
