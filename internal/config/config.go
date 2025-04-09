package config

import (
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	Topic        string
	GroupID      string
}

func LoadConfig() *Config {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	return &Config{
		KafkaBrokers: strings.Split(brokers, ","),
		Topic:        "test-topic",
		GroupID:      "test-group",
	}
}
