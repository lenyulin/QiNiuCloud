package producer

type AddEvent struct {
	Partition int32
	Offset    string
	TimeStamp int64
	Topic     string
	Retry     int64
	DATA      interface{}
}

const (
	TopicInsertModelInfoToDBEvent = "topic_insert_model_info_event"
)
