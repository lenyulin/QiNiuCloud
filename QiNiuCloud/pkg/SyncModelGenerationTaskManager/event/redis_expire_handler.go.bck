package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type RedisExpireHandler struct {
	producer sarama.SyncProducer
	redis    *redis.Client
	mu       sync.Mutex
}
type ExpireEvent struct {
	Key       string    `json:"key"`
	EventTime time.Time `json:"event_time"`
}

func (e *RedisExpireHandler) NewRedisExpireHandler(producer sarama.SyncProducer, redis *redis.Client) *RedisExpireHandler {
	return &RedisExpireHandler{
		producer: producer,
		redis:    redis,
	}
}
func (e *RedisExpireHandler) subscribeExpireEvents(ctx context.Context) {
	pubsub := e.redis.Subscribe(ctx, "__keyevent@0__:expired")
	defer pubsub.Close()
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
		return
	}
	ch := pubsub.Channel()
	fmt.Println("开始监听过期事件...")
	for msg := range ch {
		event := ExpireEvent{
			Key:       msg.Payload,
			EventTime: time.Now(),
		}
		// 序列化事件
		data, err := e.redis.Get(context.Background(), msg.Payload).Result()
		if err != nil {
			fmt.Println(err)
			continue
		}
		tccEvt := &AddTCCEvent{}
		err = json.Unmarshal([]byte(data), &tccEvt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tccEvt.Status = StatusConfirming
		bytes, _ := json.Marshal(tccEvt)
		// 写入Kafka（带重试）
		_, _, err = e.producer.SendMessage(&sarama.ProducerMessage{
			Topic: TopicWatchEvent,
			Key:   sarama.StringEncoder(msg.Payload),
			Value: sarama.ByteEncoder(bytes),
		})
		if err != nil {
			fmt.Printf("Kafka写入失败，降级写入本地文件: %v\n", err)
			// 降级：写入本地文件（防止彻底丢失）
			//saveToLocalFile(event.Key, eventData)
		} else {
			fmt.Printf("事件已持久化: %s\n", event.Key)
		}
	}
}

// 工具函数：重试逻辑
func (e *RedisExpireHandler) retry(maxRetries int, interval time.Duration, fn func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < maxRetries-1 {
			time.Sleep(interval)
		}
	}
	return err
}

// 定期补偿检查：扫描可能未被处理的过期键
func (e *RedisExpireHandler) compensateCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		keys, err := e.redis.Keys(ctx, "tcc:*").Result()
		if err != nil {
			fmt.Printf("补偿扫描失败: %v\n", err)
			continue
		}
		for _, key := range keys {
			// 检查键是否已过期（TTL <= 0）
			ttl, err := e.redis.TTL(ctx, key).Result()
			if err != nil || ttl > 0 {
				continue
			}
			fmt.Printf("补偿处理过期键: %s\n", key)
		}
	}
}
