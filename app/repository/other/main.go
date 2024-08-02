package other_repository

import (
	"workspace-server/app/repository"
)

type OtherRepositoryCollections struct {
	RedisRepo repository.RedisRepository
}

func RegisterOtherRepositories() OtherRepositoryCollections {
	return OtherRepositoryCollections{
		RedisRepo: NewRedisRepository(),
	}
}
