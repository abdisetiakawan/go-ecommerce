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

func NewKafkaProducer(viper *viper.Viper) (*KafkaProducer, error) {
    kafkaConfig := kafka.NewKafkaConfig(viper)
    producer, err := kafka.NewKafkaProducer(kafkaConfig)
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