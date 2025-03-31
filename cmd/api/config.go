package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	addr string
}

func Load() *config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	return &config{
		addr: os.Getenv("Addr"),
	}
}
