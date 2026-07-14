package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Application ApplicationConfig `mapstructure:"application"`

	Kafka KafkaConfig `mapstructure:"kafka"`
	Redis RedisConfig `mapstructure:"redis"`
}

type ApplicationConfig struct {
	Name string `mapstructure:"name"`
}

// services configuration

type KafkaConfig struct {
	Seeds []string   `mapstructure:"seeds"`
	Topic KafkaTopic `mapstructure:"publishTopics"`
}

type KafkaTopic struct {
	Route       string `mapstructure:"route"`
	Variant     string `mapstructure:"variant"`
	Stop        string `mapstructure:"stop"`
	Geolocation string `mapstructure:"geolocation"`
}

type RedisConfig struct {
	Address string `mapstructure:"address"`
}

func SetDefault() {
	viper.SetDefault("application.name", "transit-radar")

	viper.SetDefault("kafka.seeds", []string{"localhost:19092"})
	viper.SetDefault("kafka.publishTopics.route", "processor.v1beta1.route")
	viper.SetDefault("kafka.publishTopics.variant", "processor.v1beta1.variant")
	viper.SetDefault("kafka.publishTopics.stop", "processor.v1beta1.stop")
	viper.SetDefault("kafka.publishTopics.geolocation", "processor.v1beta1.geolocation")

	viper.SetDefault("redis.address", "localhost:6379")

	viper.SetDefault("task.data.enable", true)
	viper.SetDefault("task.data.crontab", "5 23,0,1,11,12,13 * * *")

	viper.SetDefault("task.geolocation.enable", true)
	viper.SetDefault("task.geolocation.crontab", "@every 30s")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}

func LoadConfig() (Config, error) {
	SetDefault()

	err := viper.ReadInConfig()
	if err != nil {
		slog.Info("no config found, using defaults...", "error", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse config: %e", err)
	}

	return config, nil
}
