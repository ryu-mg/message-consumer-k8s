package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/ryu-mg/message-consumer-k8s/internal/config"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	config   *config.Config
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		consumer: consumerGroup,
		config:   cfg,
	}, nil
}

func (c *Consumer) Start() error {
	handler := &ConsumerGroupHandler{}
	ctx := context.Background()

	log.Println("Consumer started")

	for {
		if err := c.consumer.Consume(ctx, []string{c.config.Topic}, handler); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s",
			string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
