package helper

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/abdisetiakawan/go-ecommerce/internal/kafka"
	"github.com/spf13/viper"
)

type KafkaProducer struct {
    producer sarama.SyncProducer
}

func NewKafkaProducer(v *viper.Viper) (*KafkaProducer, error) {
    config := kafka.NewKafkaConfig(v)
    
    saramaConfig, err := kafka.NewSaramaConfig(config)
    if err != nil {
        return nil, err
    }
    
    // Producer specific configs
    saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
    saramaConfig.Producer.Return.Successes = true
    saramaConfig.Producer.Return.Errors = true
    
    producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
    if err != nil {
        return nil, err
    }
    
    return &KafkaProducer{producer: producer}, nil
}

func (k *KafkaProducer) SendMessage(ctx context.Context, message interface{}, topic string) error {
    jsonMessage, err := json.Marshal(message)
    if err != nil {
        return err
    }

    _, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.ByteEncoder(jsonMessage),
    })
    return err
}

func (k *KafkaProducer) Close() error {
    if k.producer != nil {
        return k.producer.Close()
    }
    return nil
}