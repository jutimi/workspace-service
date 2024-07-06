package postgres_repository

import (
	"oauth-server/app/repository"

	"gorm.io/gorm"
)

type PostgresRepositoryCollections struct {
	PostgresOAuthRepo repository.OAuthRepository
}

func RegisterPostgresRepositories(db *gorm.DB) PostgresRepositoryCollections {
	postgresOAuthRepo := NewPostgresOAuthRepository(db)

	return PostgresRepositoryCollections{
		PostgresOAuthRepo: postgresOAuthRepo,
	}
}
