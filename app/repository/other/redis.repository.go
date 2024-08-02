package other_repository

import (
	"context"

	"github.com/jutimi/workspace-server/app/repository"
)

type redisRepository struct {
}

func NewRedisRepository() repository.RedisRepository {
	return &redisRepository{}
}

func (r *redisRepository) Set(ctx context.Context, key string, value string) error {
	return nil
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	return nil
}
