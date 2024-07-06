package postgres_repository

import (
	"gorm.io/gorm"
)

type PostgresRepositoryCollections struct {
}

func RegisterPostgresRepositories(db *gorm.DB) PostgresRepositoryCollections {

	return PostgresRepositoryCollections{}
}
