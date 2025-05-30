package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr      string
	DB        dbconfig
	KafkaAddr string
	RateLimit ratelimitConfig
}

type dbconfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
}

type ratelimitConfig struct {
	TimeFrame      time.Duration
	RequestPerTime int
	Enabled        bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	return &Config{
		KafkaAddr: os.Getenv("kafkaAddr"),
		Addr:      os.Getenv("addr"),
		DB: dbconfig{
			Addr:         os.Getenv("DBaddr"),
			MaxOpenConns: 30,
			MaxIdleConns: 30,
		},
		RateLimit: ratelimitConfig{
			Enabled:        true,
			RequestPerTime: 20,
			TimeFrame:      time.Second * 5,
		},
	}
}
