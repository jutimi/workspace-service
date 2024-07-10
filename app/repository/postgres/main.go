package postgres_repository

import (
	"workspace-server/app/repository"

	"gorm.io/gorm"
)

type PostgresRepositoryCollections struct {
	WorkspaceRepo           repository.WorkspaceRepository
	UserWorkspaceRepo       repository.UserWorkspaceRepository
	UserWorkspaceDetailRepo repository.UserWorkspaceDetailRepository
}

func RegisterPostgresRepositories(db *gorm.DB) PostgresRepositoryCollections {
	return PostgresRepositoryCollections{
		WorkspaceRepo:           NewWorkspaceRepository(db),
		UserWorkspaceRepo:       NewUserWorkspaceRepository(db),
		UserWorkspaceDetailRepo: NewUserWorkspaceDetailRepository(db),
	}
}
