package clients

import (
	"errors"
	"fmt"

	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/store"
	"github.com/IBM/sarama"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Clients struct {
	Kafka sarama.Client
	Redis *redis.Client

	Asynq *asynq.Client
	Store store.Store
}

func InitClients(config *config.WorkerConfig) (*Clients, error) {
	var clients Clients

	if err := clients.initKafka(&config.Kafka); err != nil {
		return nil, fmt.Errorf("cannot init sarama client: %w", err)
	}

	if err := clients.initRedis(&config.Redis); err != nil {
		return nil, fmt.Errorf("cannot init redis client: %w", err)
	}

	if err := clients.initAsynq(&config.Redis); err != nil {
		return nil, fmt.Errorf("cannot init asynq client: %w", err)
	}

	return &clients, nil
}

func (c *Clients) Close() error {
	var err error

	if c.Asynq != nil {
		err = errors.Join(err, c.Asynq.Close())
	}

	if c.Kafka != nil {
		err = errors.Join(err, c.Kafka.Close())
	}

	if c.Redis != nil {
		err = errors.Join(err, c.Redis.Close())
	}

	return err
}

func (c *Clients) initKafka(config *config.KafkaConfig) error {
	var err error

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true

	c.Kafka, err = sarama.NewClient(config.Seeds, saramaConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Clients) initRedis(config *config.RedisConfig) error {
	c.Redis = redis.NewClient(&redis.Options{
		Addr: config.Address,
	})

	c.Store = store.NewRedisStore(c.Redis)
	return nil
}

func (c *Clients) initAsynq(config *config.RedisConfig) error {
	c.Asynq = asynq.NewClient(asynq.RedisClientOpt{
		Addr: config.Address,
	})

	return nil
}
