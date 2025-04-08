package config

type Config struct {
	KafkaBrokers []string
	Topic        string
	GroupID      string
}

func LoadConfig() *Config {
	return &Config{
		KafkaBrokers: []string{"localhost:9092"}, // 기본값
		Topic:        "test-topic",
		GroupID:      "test-group",
	}
}
