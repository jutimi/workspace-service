package postgres_repository

import (
	"github.com/jutimi/workspace-server/app/repository"

	"gorm.io/gorm"
)

type PostgresRepositoryCollections struct {
	WorkspaceRepo                 repository.WorkspaceRepository
	UserWorkspaceRepo             repository.UserWorkspaceRepository
	UserWorkspaceDetailRepo       repository.UserWorkspaceDetailRepository
	OrganizationRepo              repository.OrganizationRepository
	UserWorkspaceOrganizationRepo repository.UserWorkspaceOrganizationRepository
}

func RegisterPostgresRepositories(db *gorm.DB) PostgresRepositoryCollections {
	return PostgresRepositoryCollections{
		WorkspaceRepo:                 NewWorkspaceRepository(db),
		UserWorkspaceRepo:             NewUserWorkspaceRepository(db),
		UserWorkspaceDetailRepo:       NewUserWorkspaceDetailRepository(db),
		OrganizationRepo:              NewOrganizationRepository(db),
		UserWorkspaceOrganizationRepo: NewUserWorkspaceOrganizationRepository(db),
	}
}
