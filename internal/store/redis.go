package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) Store {
	return &redisStore{client}
}

func (s *redisStore) Set(ctx context.Context, key string, value any) error {
	if err := s.client.HSet(ctx, key, value).Err(); err != nil {
		return err
	}

	return nil
}

func (s *redisStore) Get(ctx context.Context, key string, value any) error {
	if err := s.client.HGetAll(ctx, key).Scan(value); err != nil {
		return err
	}

	return nil
}

func (s *redisStore) Add(ctx context.Context, set string, value ...string) error {
	if err := s.client.SAdd(ctx, set, value).Err(); err != nil {
		return err
	}

	return nil
}

func (s *redisStore) IsMember(ctx context.Context, set string, value string) (bool, error) {
	cmd := s.client.SIsMember(ctx, set, value)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val(), nil
}

func (s *redisStore) Members(ctx context.Context, set string) ([]string, error) {
	cmd := s.client.SMembers(ctx, set)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	return cmd.Val(), nil
}
