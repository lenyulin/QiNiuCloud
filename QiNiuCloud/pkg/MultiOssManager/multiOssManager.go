package MultiOssManager

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	oss "QiNiuCloud/QiNiuCloud/pkg/ossx"
	"context"
	"errors"
	"sync/atomic"
	"time"
)

type MultiOssManager interface {
	Upload(filename string)
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
	mainPool             []*OSSObject
	mainPoolSize         int
	backupPool           []*OSSObject
	backupPoolSize       int
	failureMainOssPool   map[int]*OSSObject
	failureBackupOssPool map[int]*OSSObject
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
	op        int
	fileName  string
	status    Status
}
type Status struct {
	IsAvailable bool
	IsFailure   bool
}
type ossManager struct {
	l       logger.ZapLogger
	OSSPool *OSSPool
	msgChan chan *ActorMsg
}

func (o *ossManager) runActor() {
	for msg := range o.msgChan {
		switch msg.op {
		case UploadFileToOSS:
			o.Upload(msg.fileName)
		case DeleteFileFromOSS:
			o.Delete()
		case FindFileFromOSS:
			o.Find()
		case ChangeOssStatusToSuccess:

		case ChangeStatusToAvailable:
		case ChangeOssStatusToFail:
		default:
			panic("unhandled default case")
		}
	}
}

func (o *ossManager) markAsFailed(object *OSSObject, version int64) {
	object.FailureCount += 1
	object.LastFailureAt = time.Now()
	object.LastModified = time.Now()
}
func (o *ossManager) markAsAvailable(object *OSSObject, version int64) {
	object.IsAvailable.Swap(true)
	object.LastModified = time.Now()
}
func (o *ossManager) getAvailableOss() (*OSSObject, error) {
	for _, obj := range o.OSSPool.mainPool {
		if obj.IsAvailable.Load() {
			return obj, nil
		}
	}
	for _, obj := range o.OSSPool.backupPool {
		if obj.IsAvailable.Load() {
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
		l:       l,
		OSSPool: ossPool,
		msgChan: ch,
	}
}

const MaxRetryCount = 3

func (o *ossManager) Upload(filename string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	availableOss, err := o.getAvailableOss()
	if err != nil {
		//处理？
	}
	go func(ctx context.Context, oss *OSSObject) {
		retry := 0
		_, _, er := oss.Client.Upload(ctx, filename)
		for er != nil && retry <= MaxRetryCount {
			//进一步处理?
			retry += 1
			time.Sleep(time.Millisecond * 30)
			_, _, er = oss.Client.Upload(ctx, filename)
		}
		if er != nil {
			if errors.Is(err, ErrOssClientFailed) {
				o.msgChan <- &ActorMsg{
					op:        ChangeOssStatusToFail,
					ossObject: oss,
					fileName:  filename,
					status: Status{
						IsAvailable: false,
						IsFailure:   true,
					},
				}
			}
			o.l.Error(err.Error())
		}
		oss.IsAvailable.Store(true)
	}(ctx, availableOss)
}

func (o *ossManager) Delete() {
	//TODO implement me
	panic("implement me")
}

func (o *ossManager) Find() {
	//TODO implement me
	panic("implement me")
}
