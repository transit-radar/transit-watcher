package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type KafkaConfig struct {
	Seeds []string `mapstructure:"seeds"`
}

func SetDefault() {
	viper.SetDefault("kafka.seeds", []string{"localhost:9092"})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}

func LoadConfig(ctx context.Context) (Config, error) {
	SetDefault()

	err := viper.ReadInConfig()
	if err != nil {
		slog.InfoContext(ctx, "no config found, using defaults...", "error", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse config: %e", err)
	}

	return config, nil
}
