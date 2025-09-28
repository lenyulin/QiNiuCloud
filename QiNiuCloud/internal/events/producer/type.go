package producer

type InteractiveEvent struct {
	Partition int32
	Offset    string
	TimeStamp int64
	Topic     string
	Retry     int64
	DATA      interface{}
}

const (
	TopicInteractiveEvent = "topic_interactive_event"
)
