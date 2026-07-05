package store

import (
	"context"
)

type Store interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string, value any) error

	Add(ctx context.Context, set string, value ...string) error
	IsMember(ctx context.Context, set string, value string) (bool, error)
	Members(ctx context.Context, set string) ([]string, error)
}
