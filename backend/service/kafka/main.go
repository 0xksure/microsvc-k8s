package kafka

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type BountyKafkaClient struct {
	Logger *zerolog.Logger
	Reader *kafka.Reader
}

type KafkaMessage struct {
	Topic string
	Key   string
	Msg   []byte
}

func (b *BountyKafkaClient) createKafkaConsumer(topic string) {
	mechanism, err := scram.Mechanism(scram.SHA512, "user1", os.Getenv("KAFKA_PASSWORD"))
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092",
			"kafka-controller-1.kafka-controller-headless.default.svc.cluster.local:9092",
			"kafka-controller-2.kafka-controller-headless.default.svc.cluster.local:9092"},
		Topic:     topic,
		MaxBytes:  10e6, // 10MB
		Logger:    b.Logger,
		Partition: 0,
		Dialer:    dialer,
	})
	defer r.Close()
	b.Reader = r
}

// GenerateKafkaConsumer generates a kafka consumer for the given topic
func (b *BountyKafkaClient) GenerateKafkaConsumer(ctx context.Context, topic string, kf chan KafkaMessage) {
	mechanism, err := scram.Mechanism(scram.SHA512, "user1", os.Getenv("KAFKA_PASSWORD"))
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092",
			"kafka-controller-1.kafka-controller-headless.default.svc.cluster.local:9092",
			"kafka-controller-2.kafka-controller-headless.default.svc.cluster.local:9092"},
		Topic:     topic,
		MaxBytes:  10e6, // 10MB
		Logger:    b.Logger,
		Partition: 0,
		Dialer:    dialer,
	})
	defer r.Close()

	b.Logger.Info().Msgf("Starting kafka consumer for topic %s", topic)
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			break
		}
		b.Logger.Info().Msgf("Received message from kafka: %s", m.Value)
		kf <- KafkaMessage{
			Topic: m.Topic,
			Key:   string(m.Key),
			Msg:   m.Value,
		}
	}
}
