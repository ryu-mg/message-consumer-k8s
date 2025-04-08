package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryu-mg/message-consumer-k8s/internal/config"
	"github.com/ryu-mg/message-consumer-k8s/internal/consumer"
)

func main() {
	cfg := config.LoadConfig()

	kafkaConsumer, err := consumer.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// Graceful shutdown을 위한 채널 설정
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// 컨슈머 실행
	go func() {
		if err := kafkaConsumer.Start(); err != nil {
			log.Printf("Error running consumer: %v", err)
		}
	}()

	// 시그널 대기
	<-sigterm
	log.Println("Terminating: via signal")

	// 컨슈머 종료
	if err := kafkaConsumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}
