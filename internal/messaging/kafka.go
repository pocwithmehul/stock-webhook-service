package messaging

import (
	"github.com/pocwithmehul/common-go-lib"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(cfg *commonlib.Config) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.Hash{},
	})
}
