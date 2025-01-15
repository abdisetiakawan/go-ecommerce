package config

import "github.com/spf13/viper"

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		Charset  string
	}
}

func LoadConfig(viper *viper.Viper) (*Config, error) {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}