package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/retry"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brooker []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr: kafka.TCP(brooker...),
		},
	}
}

func (k *KafkaService) Produce(topic string, key, value []byte) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.Brokers,
		Topic:   topic,
	})
	defer writer.Close()

	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
}

// cannot directly create topic in kafka we sure here topic is existed(creating it)
func (k *KafkaService) EnsureTopicExists(topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var conn *kafka.Conn
	retry.Retry(ctx, func() error {
		var err error
		conn, err = kafka.Dial("tcp", "kafka:9092")
		if err != nil {
			return err
		}
		return nil
	})

	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
}
