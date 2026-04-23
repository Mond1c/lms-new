package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Config struct {
	Brokers       string
	ClientID      string
	ConsumerGroup string
	Topics        []string
}

func newKafkaOpts(brokers, clientID string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(strings.Split(brokers, ",")...),
		kgo.ClientID(clientID),
	}
}

func NewProducer(cfg Config) (*kgo.Client, error) {
	const batchingTimeout = 5 * time.Millisecond

	opts := newKafkaOpts(cfg.Brokers, cfg.ClientID)
	opts = append(
		opts,
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerLinger(batchingTimeout),
	)

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("kafka producer: %w", err)
	}
	return client, nil
}

func NewConsumer(cfg Config) (*kgo.Client, error) {
	if cfg.ConsumerGroup == "" {
		return nil, fmt.Errorf("kafka consumer group required")
	}
	if len(cfg.Topics) == 0 {
		return nil, fmt.Errorf("kafka consumer topics required")
	}
	opts := newKafkaOpts(cfg.Brokers, cfg.ClientID)
	opts = append(
		opts,
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.DisableAutoCommit(),
	)

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("kafka consumer: %w", err)
	}
	return client, nil
}

func Ping(ctx context.Context, client *kgo.Client) error {
	return client.Ping(ctx)
}
