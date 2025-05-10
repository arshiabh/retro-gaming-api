package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

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

func (k *KafkaService) StartConsumer(ctx context.Context, reader *kafka.Reader) {
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down kafka service")
			return
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("error reading message")
				log.Println(err)
			}
			log.Printf("message recevied: %s\n", string(m.Value))
		}
	}
}
