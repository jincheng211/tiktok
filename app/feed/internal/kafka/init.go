package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

var (
	KafkaConn *kafka.Conn
)

func InitKafka(broker, topic string) error {

	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
	if err != nil {
		return err
	}
	KafkaConn = conn
	// 在这里进行其他Kafka初始化操作
	fmt.Println("成功初始化Kafka连接")

	return nil
}
