package helper

import (
	postgres_repository "workspace-server/app/repository/postgres"
)

type HelperCollections struct {
}

func RegisterHelpers(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) HelperCollections {

	return HelperCollections{}
}
