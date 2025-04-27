package kafka

import (
	"context"

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

func (k *Client) CreateReader(topic string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.Brokers,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return reader
}
