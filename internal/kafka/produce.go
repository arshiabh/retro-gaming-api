package kafka

import (
	"context"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Produce(string, []byte, []byte) error
	EnsureTopicExists(string) error
	CreateReader(string, string) *kafka.Reader
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
	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		return err
	}
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
