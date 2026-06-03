package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"digit-oss/noc-services/internal/domain"

	"github.com/Shopify/sarama"
)

// KafkaProducer implements domain.Producer using a sarama SyncProducer.
// Translates Producer.java → kafkaTemplate.send(topic, value)
type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// Compile-time interface check.
var _ domain.Producer = (*KafkaProducer)(nil)

// Push serializes the value to JSON and publishes it to the given Kafka topic.
func (p *KafkaProducer) Push(topic string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal kafka message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send to %s: %w", topic, err)
	}

	log.Printf("[KafkaProducer] topic=%s partition=%d offset=%d", topic, partition, offset)
	return nil
}

// NewKafkaProducer creates a new sarama SyncProducer connected to the given brokers.
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &KafkaProducer{Producer: producer}, nil
}

// Close shuts down the producer gracefully.
func (p *KafkaProducer) Close() error {
	return p.Producer.Close()
}
