package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
	DB   dbconfig
}

type dbconfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	return &Config{
		Addr: os.Getenv("addr"),
		DB: dbconfig{
			Addr:         os.Getenv("DBaddr"),
			MaxOpenConns: 30,
			MaxIdleConns: 30,
		},
	}
}
