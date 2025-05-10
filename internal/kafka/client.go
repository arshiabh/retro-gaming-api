package kafka

type KafkaService struct {
	Brokers []string
}

func NewKafkaService(brokers []string) *KafkaService {
	return &KafkaService{
		Brokers: brokers,
	}
}
