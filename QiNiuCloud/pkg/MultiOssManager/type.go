package MultiOssManager

import (
	oss "QiNiuCloud/QiNiuCloud/pkg/ossx"
	"context"
	"time"
)

type MultiOssManager interface {
	Upload(ctx context.Context, uid int64, filename string) error
	Delete(ctx context.Context, uid int64) error
	Find(ctx context.Context, uid int64) (string, error)
}
type OSSObject struct {
	Client        *oss.OSSHandler // OSS客户端
	LastSuccess   bool            // 上次访问是否成功
	FailureCount  int             // 失效次数
	RetryCount    int             // 重试次数
	LastFailureAt time.Time       // 最近失效时间
	LastModified  time.Time       // 最近修改时间
	Version       int64           // 乐观锁版本号，用于并发控制
}
type OSSPool struct {
	mainPool       []*OSSObject
	mainPoolSize   int
	backupPool     []*OSSObject
	backupPoolSize int
	retryPool      map[string]*OSSObject
}

type ActorMsg struct {
	obj       *OSSObject
	version   int64
	isSucceed bool
}
type ossManager struct {
	OSSPool *OSSPool
	msgChan chan *ActorMsg
}

func (o *ossManager) runActor() {
	for msg := range o.msgChan {
		switch msg.isSucceed {
		case false:
			o.markAsFailed(msg.obj, msg.version)
		case true:
			o.markAsAvailable(msg.obj, msg.version)
		}
	}
}

func (o *ossManager) markAsFailed(object *OSSObject, version int64) {
	if object.Version < version {
		return
	}
	object.FailureCount += 1
	object.LastFailureAt = time.Now()
	object.LastModified = time.Now()
	object.Version = version + 1
	object.Client
}
func (o *ossManager) markAsAvailable(object *OSSObject, version int64) {

}

func NewOSSManager(mainOss []*oss.OSSHandler, backupOss []*oss.OSSHandler) MultiOssManager {
	ch := make(chan *ActorMsg, 1000)
	ossPool := &OSSPool{
		mainPool:   make([]*OSSObject, 0),
		backupPool: make([]*OSSObject, 0),
	}
	for _, o := range mainOss {
		if o != nil {
			ossObj := &OSSObject{
				Client:        o,
				LastSuccess:   false,
				FailureCount:  0,
				RetryCount:    0,
				LastFailureAt: time.Time{},
				LastModified:  time.Now(),
				Version:       0,
			}
			ossPool.mainPool = append(ossPool.mainPool, ossObj)
			ossPool.mainPoolSize += 1
		}
	}
	for _, o := range backupOss {
		if o != nil {
			ossObj := &OSSObject{
				Client:        o,
				LastSuccess:   false,
				FailureCount:  0,
				RetryCount:    0,
				LastFailureAt: time.Time{},
				LastModified:  time.Now(),
				Version:       0,
			}
			ossPool.backupPool = append(ossPool.backupPool, ossObj)
			ossPool.backupPoolSize += 1
		}
	}
	return &ossManager{
		OSSPool: ossPool,
		msgChan: ch,
	}
}
func (o *ossManager) Upload(ctx context.Context, uid int64, filename string) error {
	//TODO implement me
	panic("implement me")
}

func (o *ossManager) Delete(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *ossManager) Find(ctx context.Context, uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}
