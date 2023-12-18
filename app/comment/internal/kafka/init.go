package kafka

import (
	"context"
	"douyin/config"
	"github.com/segmentio/kafka-go"
	"sync"
)

type Kafka struct {
	reader         *kafka.Reader
	cond           *sync.Cond
	messageHandled bool
}

var (
	KafkaCli *Kafka
)

// 初始化 KafkaManager
func InitKafka() {
	// 执行初始化逻辑
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{config.Conf.Kafka.Addr},
		Topic:     "CommentActionMessage",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	KafkaCli = &Kafka{
		reader: r,
		cond:   sync.NewCond(&sync.Mutex{}),
	}
}

const (
	NoMessage   = false
	HaveMessage = true
)

func InitConsumer() {
	go KafkaCli.Consumer()
}

func (k *Kafka) Consumer() {
	for {
		k.cond.L.Lock()
		if k.messageHandled == NoMessage {
			// 没有消息卡柱
			k.cond.Wait()
		}
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		// 在这里处理消息
		_ = m.Topic

		k.cond.Signal()
		k.cond.L.Unlock()
	}
}
