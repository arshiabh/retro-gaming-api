package kafka

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/retry"
	"github.com/segmentio/kafka-go"
)

func (k *KafkaService) Produce(ctx context.Context, topic string, key, value string) error {
	k.mu.Lock()
	writer, exists := k.Producermap[topic]
	if !exists {
		writer = &kafka.Writer{
			Addr:     kafka.TCP(k.Brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
		k.Producermap[topic] = writer
	}
	k.mu.Unlock()

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	})
}

func SendAsync(wg *sync.WaitGroup, topic, key, value string, sender *KafkaService) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := sender.Produce(ctx, topic, key, value); err != nil {
			log.Println(err)
		}
	}()
}

func (k *KafkaService) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()
	for topic, writer := range k.Producermap {
		if err := writer.Close(); err != nil {
			log.Printf("Error closing writer for topic %s: %v", topic, err)
		}
	}
}

// cannot directly create topic in kafka we sure here topic is existed(creating it)
func (k *KafkaService) EnsureTopicExists(topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var conn *kafka.Conn
	retry.Retry(ctx, func() error {
		var err error
		conn, err = kafka.Dial("tcp", "localhost:9092")
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
