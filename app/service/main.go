package service

import (
	"workspace-server/app/helper"
	postgres_repository "workspace-server/app/repository/postgres"
)

type ServiceCollections struct {
	OAuthSvc OAuthService
}

func RegisterServices(
	helpers helper.HelperCollections,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
) ServiceCollections {
	oauthSvc := NewOAuthService(helpers, postgresRepo)

	return ServiceCollections{
		OAuthSvc: oauthSvc,
	}
}
