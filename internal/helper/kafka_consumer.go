package helper

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/abdisetiakawan/go-ecommerce/internal/kafka"
	"github.com/spf13/viper"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
}

func NewKafkaConsumer(viper *viper.Viper) (*KafkaConsumer, error) {
	kafkaConfig := kafka.NewKafkaConfig(viper)
	consumer, err := sarama.NewConsumer(kafkaConfig.Brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer}, nil
}

func (k *KafkaConsumer) Consume(ctx context.Context, topic string) (sarama.PartitionConsumer, error) {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	return partitionConsumer, nil
}