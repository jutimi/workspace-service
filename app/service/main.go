package service

import (
	"workspace-server/app/helper"
	other_repository "workspace-server/app/repository/other"
	postgres_repository "workspace-server/app/repository/postgres"
)

type ServiceCollections struct {
	WorkspaceSvc     WorkspaceService
	UserWorkspaceSvc UserWorkspaceService
}

func RegisterServices(
	helpers helper.HelperCollections,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
	otherRepo other_repository.OtherRepositoryCollections,
) ServiceCollections {
	return ServiceCollections{
		WorkspaceSvc:     NewWorkspaceService(helpers, postgresRepo),
		UserWorkspaceSvc: NewUserWorkspaceService(helpers, postgresRepo),
	}
}
