package kafka

import (
	"context"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	Brokers []string
}

func NewClient(brokers []string) *Client {
	return &Client{
		Brokers: brokers,
	}
}

func (k *Client) Produce(topic string, key, value []byte) error {
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

func (k *Client) CreateReader(groupID, topic string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.Brokers,
		Topic:   topic,
		GroupID: groupID,
		//start reading new message
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})

	return reader
}

// cannot directly create topic in kafka we sure here topic is existed(creating it)
func (k *Client) EnsureTopicExists(topic string) error {
	conn, err := kafka.Dial("tcp", "localhost:9092")
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
