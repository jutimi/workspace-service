package service

import (
	"workspace-server/app/helper"
	postgres_repository "workspace-server/app/repository/postgres"
)

type ServiceCollections struct {
	WorkspaceSvc WorkspaceService
}

func RegisterServices(
	helpers helper.HelperCollections,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
) ServiceCollections {
	return ServiceCollections{
		WorkspaceSvc: NewWorkspaceService(helpers, postgresRepo),
	}
}
