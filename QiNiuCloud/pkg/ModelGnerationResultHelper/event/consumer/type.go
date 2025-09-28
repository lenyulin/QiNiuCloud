package consumer

type Consumer interface {
	Start() error
}

type AddEvent struct {
	Partition int32
	Offset    string
	TimeStamp int64
	Topic     string
	Retry     int64
	DATA      interface{}
}

const (
	TopicWatchEvent = "model_generate_result_watch"
)
