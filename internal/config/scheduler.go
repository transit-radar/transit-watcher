package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type SchedulerConfig struct {
	Application ApplicationConfig `mapstructure:"application"`

	Redis RedisConfig `mapstructure:"redis"`

	Task TaskConfig `mapstructure:"task"`
}

func NewSchedulerConfig() *SchedulerConfig {
	c := &SchedulerConfig{}
	c.SetDefault()
	return c
}

func (c *SchedulerConfig) SetDefault() {
	viper.SetDefault("application.name", "app.transitradar.watcher.scheduler")

	// redis
	viper.SetDefault("redis.address", "localhost:6379")

	// tasks
	viper.SetDefault("task.data.enable", true)
	viper.SetDefault("task.data.crontab", "5 23,0,1,11,12,13 * * *")

	viper.SetDefault("task.multiGoGeolocation.enable", true)
	viper.SetDefault("task.multiGoGeolocation.crontab", "@every 30s")

	viper.SetDefault("task.ttgtGeolocation.enable", true)
	viper.SetDefault("task.ttgtGeolocation.crontab", "@every 1s")

	viper.SetConfigName("scheduler")
	viper.AddConfigPath("./config")
}

func (c *SchedulerConfig) Load() error {
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
