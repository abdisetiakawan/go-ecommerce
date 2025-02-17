package kafka

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
    Brokers []string
}

func NewKafkaConfig(viper *viper.Viper) *KafkaConfig {
    return &KafkaConfig{
        Brokers: viper.GetStringSlice("KAFKA_BROKERS"),
    }
}

func NewKafkaProducer(config *KafkaConfig) (sarama.SyncProducer, error) {
    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.Return.Successes = true
    saramaConfig.Producer.Return.Errors = true

    producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
    if err != nil {
        return nil, err
    }

    return producer, nil
}