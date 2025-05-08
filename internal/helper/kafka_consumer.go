package helper

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/abdisetiakawan/go-ecommerce/internal/kafka"
	"github.com/spf13/viper"
)

type KafkaConsumer struct {
	consumer    sarama.ConsumerGroup
	ready       chan bool
	consumerID  string
	config      *kafka.KafkaConnectionConfig
}

func NewKafkaConsumer(v *viper.Viper, consumerID string) (*KafkaConsumer, error) {
	config := kafka.NewKafkaConnectionConfig(v)
	
	kafkaConfig := &kafka.KafkaConfig{
		Brokers:       config.Brokers,
		ConsumerGroup: config.ConsumerGroup,
	}
	saramaConfig, err := kafka.NewSaramaConfig(kafkaConfig)
	if err != nil {
		return nil, err
	}
	
	// Consumer specific configs
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	
	// Use unique consumer group ID for each consumer
	groupID := fmt.Sprintf("%s-%s", config.ConsumerGroup, consumerID)
	group, err := sarama.NewConsumerGroup(config.Brokers, groupID, saramaConfig)
	if err != nil {
		return nil, err
	}
	
	return &KafkaConsumer{
		consumer:   group,
		ready:      make(chan bool),
		consumerID: consumerID,
		config:     config,
	}, nil
}

func (k *KafkaConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	for {
		err := k.consumer.Consume(ctx, topics, handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		k.ready = make(chan bool)
	}
}

func (k *KafkaConsumer) Close() error {
    if k.consumer != nil {
        return k.consumer.Close()
    }
    return nil
}