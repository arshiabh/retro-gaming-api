package kafka

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/retry"
	"github.com/segmentio/kafka-go"
)

var (
	Producermap = make(map[string]*kafka.Writer)
	mu          sync.Mutex
)


func (k *KafkaService) Produce(topic string, key, value string) error {
	mu.Lock()
	writer, exists := Producermap[topic]
	defer writer.Close()
	if !exists {
		writer = &kafka.Writer{
			Addr:     kafka.TCP(k.Brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
		Producermap[topic] = writer
	}
	mu.Unlock()

	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
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
