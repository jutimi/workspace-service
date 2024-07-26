package service

import (
	"context"
	"workspace-server/app/entity"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	postgres_repository "workspace-server/app/repository/postgres"
	client_grpc "workspace-server/grpc/client"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/jutimi/grpc-service/oauth"
)

type userWorkspaceService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
	clientGRPC   client_grpc.OAuthClient
}

func NewUserWorkspaceService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	clientGRPC client_grpc.OAuthClient,
) UserWorkspaceService {
	return &userWorkspaceService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
		clientGRPC:   clientGRPC,
	}
}

func (s *userWorkspaceService) CreateUserWorkspace(ctx context.Context, data *model.CreateUserWorkspaceRequest) (*model.CreateUserWorkspaceResponse, error) {
	workspacePayload, err := utils.GetScopeContext[*utils.WorkspacePayload](ctx, string(utils.WORKSPACE_CONTEXT_KEY))
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Create user account
	user, err := s.clientGRPC.CreateUser(ctx, &oauth.CreateUserParams{
		PhoneNumber: data.PhoneNumber,
		Email:       data.PhoneNumber,
		Password:    entity.DEFAULT_PASSWORD,
	})
	if err != nil {
		return nil, err
	}

	tx := database.BeginPostgresTransaction()
	userWS := entity.NewUserWorkspace()
	userWS.Role = entity.ROLE_USER
	userWS.WorkspaceID = workspacePayload.WorkspaceID
	userWS.UserID = user.ID
	userWS.Email = data.Email
	userWS.PhoneNumber = data.PhoneNumber

	return nil, nil
}

func (s *userWorkspaceService) UpdateUserWorkspace(ctx context.Context, data *model.UpdateUserWorkspaceRequest) (*model.UpdateUserWorkspaceResponse, error) {
	return nil, nil
}

func (s *userWorkspaceService) InactiveUserWorkspace(ctx context.Context, data *model.InactiveUserWorkspaceRequest) (*model.InactiveUserWorkspaceResponse, error) {
	return nil, nil
}

func (s *userWorkspaceService) RemoveUserWorkspace(ctx context.Context, data *model.RemoveUserWorkspaceRequest) (*model.RemoveUserWorkspaceResponse, error) {
	return nil, nil
}
