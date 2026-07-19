package config

import "time"

type ApplicationConfig struct {
	Name string `mapstructure:"name"`
}

// services configuration

type KafkaConfig struct {
	Seeds []string   `mapstructure:"seeds"`
	Topic KafkaTopic `mapstructure:"publishTopics"`
}

type KafkaTopic struct {
	Route        string `mapstructure:"route"`
	Variant      string `mapstructure:"variant"`
	VariantStops string `mapstructure:"variantStops"`
	Stop         string `mapstructure:"stop"`
	Geolocation  string `mapstructure:"geolocation"`
}

type RedisConfig struct {
	Address string `mapstructure:"address"`
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
