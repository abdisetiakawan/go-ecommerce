package kafka

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
    Brokers       []string
    ConsumerGroup string
    Username      string
    Password      string
    ClientID      string
    Version       string
}

func NewKafkaConfig(v *viper.Viper) *KafkaConfig {
    return &KafkaConfig{
        Brokers:       v.GetStringSlice("KAFKA_BROKERS"),
        ConsumerGroup: v.GetString("KAFKA_CONSUMER_GROUP"),
        Username:      v.GetString("KAFKA_USERNAME"),
        Password:      v.GetString("KAFKA_PASSWORD"),
        ClientID:      v.GetString("KAFKA_CLIENT_ID"),
        Version:       v.GetString("KAFKA_VERSION"),
    }
}

func NewSaramaConfig(config *KafkaConfig) (*sarama.Config, error) {
    saramaConfig := sarama.NewConfig()
    
    // Set client ID
    if config.ClientID != "" {
        saramaConfig.ClientID = config.ClientID
    }
    
    // Set SASL credentials if provided
    if config.Username != "" && config.Password != "" {
        saramaConfig.Net.SASL.Enable = true
        saramaConfig.Net.SASL.User = config.Username
        saramaConfig.Net.SASL.Password = config.Password
    }
    
    // Set Kafka version if provided
    if config.Version != "" {
        version, err := sarama.ParseKafkaVersion(config.Version)
        if err != nil {
            return nil, err
        }
        saramaConfig.Version = version
    }
    
    return saramaConfig, nil
}

func NewKafkaProducer(config *KafkaConfig) (sarama.SyncProducer, error) {
    saramaConfig, err := NewSaramaConfig(config)
    if err != nil {
        return nil, err
    }

    producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
    if err != nil {
        return nil, err
    }

    return producer, nil
}