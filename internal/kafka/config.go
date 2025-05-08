package kafka

import "github.com/spf13/viper"

type KafkaConnectionConfig struct {
	Brokers          []string
	ConsumerGroup    string
	AutoOffsetReset  string
	SecurityProtocol string
	SASLMechanism    string
	SASLUsername     string
	SASLPassword     string
	RetryMax         int
	RequestTimeoutMs int
}

func NewKafkaConnectionConfig(viper *viper.Viper) *KafkaConnectionConfig {
	return &KafkaConnectionConfig{
		Brokers:          viper.GetStringSlice("KAFKA_BROKERS"),
		ConsumerGroup:    viper.GetString("KAFKA_CONSUMER_GROUP"),
		AutoOffsetReset:  viper.GetString("KAFKA_AUTO_OFFSET_RESET"),
		SecurityProtocol: viper.GetString("KAFKA_SECURITY_PROTOCOL"),
		SASLMechanism:    viper.GetString("KAFKA_SASL_MECHANISM"),
		SASLUsername:     viper.GetString("KAFKA_SASL_USERNAME"),
		SASLPassword:     viper.GetString("KAFKA_SASL_PASSWORD"),
		RetryMax:         viper.GetInt("KAFKA_RETRY_MAX"),
		RequestTimeoutMs: viper.GetInt("KAFKA_REQUEST_TIMEOUT_MS"),
	}
}