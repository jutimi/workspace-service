package helper

import (
	other_repository "workspace-server/app/repository/other"
	postgres_repository "workspace-server/app/repository/postgres"
)

type HelperCollections struct {
	UserWSHelper       UserWSHelper
	OrganizationHelper OrganizationHelper
}

func RegisterHelpers(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	otherRepo other_repository.OtherRepositoryCollections,
) HelperCollections {
	return HelperCollections{
		UserWSHelper:       NewUserWSHelper(postgresRepo),
		OrganizationHelper: NewOrganizationHelper(postgresRepo),
	}
}
