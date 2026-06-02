package messaging

import (
	"github.com/pocwithmehul/stock-webhook-service/internal/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(cfg *config.Config) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.Hash{},
	})
}
