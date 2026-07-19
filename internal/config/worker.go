package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type WorkerConfig struct {
	Application ApplicationConfig `mapstructure:"application"`

	Kafka KafkaConfig `mapstructure:"kafka"`
	Redis RedisConfig `mapstructure:"redis"`
}

func NewWorkerConfig() *WorkerConfig {
	c := &WorkerConfig{}
	c.SetDefault()
	return c
}

func (c *WorkerConfig) SetDefault() {
	viper.SetDefault("application.name", "app.transitradar.watcher.worker")

	// kafka
	viper.SetDefault("kafka.seeds", []string{"localhost:19092"})
	viper.SetDefault("kafka.publishTopics.route", "processor.v1beta1.route")
	viper.SetDefault("kafka.publishTopics.variant", "processor.v1beta1.variant")
	viper.SetDefault("kafka.publishTopics.variantStops", "processor.v1beta1.variantstops")
	viper.SetDefault("kafka.publishTopics.stop", "processor.v1beta1.stop")
	viper.SetDefault("kafka.publishTopics.geolocation", "processor.v1beta1.geolocation")

	// redis
	viper.SetDefault("redis.address", "localhost:6379")

	viper.SetConfigName("worker")
	viper.AddConfigPath("./config")
}

func (c *WorkerConfig) Load() error {
	err := viper.ReadInConfig()
	if err != nil {
		slog.Info("no config found, using defaults...", "error", err)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("cannot parse config: %e", err)
	}

	return nil
}
