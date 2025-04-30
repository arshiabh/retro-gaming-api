package module

import (
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type Dependencies struct {
	Store *store.Storage
	Kafka *kafka.Client
}
