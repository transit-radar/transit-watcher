package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

type SchedulerConfig struct {
	Application ApplicationConfig `mapstructure:"application"`

	Redis RedisConfig `mapstructure:"redis"`

	Task TaskConfig `mapstructure:"task"`
}

type TaskConfig struct {
	// Aggregate performs routes, variants, and stops data aggregation from
	// multiple supported sources (currently GoBus and EBMS)
	Aggregate TaskSpec `mapstructure:"aggregate"`

	// Geolocation retrieve geolocation data per cached routes from MultiGo
	MultiGoGeolocation TaskSpec `mapstructure:"multiGoGeolocation"`
	TTGTGeolocation    TaskSpec `mapstructure:"ttgtGeolocation"`
}

type TaskSpec struct {
	// Whether to enable or disable the task
	Enable bool `mapstructure:"enable"`

	// Crontab spec to trigger the task periodically
	Crontab string `mapstructure:"crontab"`

	MaxRetry *int           `mapstructure:"maxRetry"`
	Unique   *time.Duration `mapstructure:"unique"`
}

func NewSchedulerConfig() *SchedulerConfig {
	c := &SchedulerConfig{}
	c.SetDefault()
	return c
}

func (c *SchedulerConfig) SetDefault() {
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
