package kafka

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	Brokers     []string
	Producermap map[string]*kafka.Writer
	mu          sync.Mutex
}

func NewKafkaService(brokers []string) *KafkaService {
	return &KafkaService{
		Brokers:     brokers,
		Producermap: make(map[string]*kafka.Writer),
		mu:          sync.Mutex{},
	}
}
