package kafka

import "github.com/segmentio/kafka-go"

func (k *KafkaService) CreateReader(groupID, topic string) *kafka.Reader {
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
