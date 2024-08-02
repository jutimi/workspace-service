package other_repository

import (
	"github.com/jutimi/workspace-server/app/repository"
)

type OtherRepositoryCollections struct {
	RedisRepo repository.RedisRepository
}

func RegisterOtherRepositories() OtherRepositoryCollections {
	return OtherRepositoryCollections{
		RedisRepo: NewRedisRepository(),
	}
}
