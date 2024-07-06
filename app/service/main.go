package service

import (
	"oauth-server/app/helper"
	postgres_repository "oauth-server/app/repository/postgres"
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
