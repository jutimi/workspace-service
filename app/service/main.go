package service

import (
	"github.com/jutimi/workspace-server/app/helper"
	other_repository "github.com/jutimi/workspace-server/app/repository/other"
	postgres_repository "github.com/jutimi/workspace-server/app/repository/postgres"
)

type ServiceCollections struct {
	WorkspaceSvc     WorkspaceService
	UserWorkspaceSvc UserWorkspaceService
	OrganizationSvc  OrganizationService
}

func RegisterServices(
	helpers helper.HelperCollections,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
	otherRepo other_repository.OtherRepositoryCollections,
) ServiceCollections {
	return ServiceCollections{
		WorkspaceSvc:     NewWorkspaceService(helpers, postgresRepo),
		UserWorkspaceSvc: NewUserWorkspaceService(helpers, postgresRepo),
		OrganizationSvc:  NewOrganizationService(helpers, postgresRepo),
	}
}
