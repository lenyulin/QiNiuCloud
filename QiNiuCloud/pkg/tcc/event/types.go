package event

type Consumer interface {
	Start() error
}

type AddTCCEvent struct {
	TCCIdx    string
	Partition int32
	Offset    string
	TimeStamp int64
	Topic     string
	Retry     int64
	Status    TransactionStatus
	DATA      interface{}
}

const (
	TopicWatchEvent     = "tcc_manger_watch"
	TCCConfirmTaskEvent = "tcc_confirm_task"
)

type TransactionStatus string

const (
	StatusTrying     TransactionStatus = "TRYING"     // 尝试中
	StatusConfirming TransactionStatus = "CONFIRMING" // 确认中
	StatusCanceling  TransactionStatus = "CANCELING"  // 取消中
	StatusCompleted  TransactionStatus = "COMPLETED"  // 已完成
	StatusCanceled   TransactionStatus = "CANCELED"   // 已取消
)
