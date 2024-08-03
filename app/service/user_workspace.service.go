package service

import (
	"context"

	"workspace-server/app/entity"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/external/client"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/jutimi/grpc-service/oauth"
)

type userWorkspaceService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserWorkspaceService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserWorkspaceService {
	return &userWorkspaceService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
	}
}

func (s *userWorkspaceService) CreateUserWorkspace(ctx context.Context, data *model.CreateUserWorkspaceRequest) (*model.CreateUserWorkspaceResponse, error) {
	workspacePayload, err := utils.GetScopeContext[*utils.WorkspacePayload](ctx, string(utils.WORKSPACE_CONTEXT_KEY))
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Create user account
	clientGRPC := client.NewOAuthClient()
	defer clientGRPC.CloseConn()
	user, err := clientGRPC.CreateUser(ctx, &oauth.CreateUserParams{
		PhoneNumber: data.PhoneNumber,
		Email:       data.PhoneNumber,
		Password:    entity.DEFAULT_PASSWORD,
	})
	if err != nil {
		return nil, err
	}
	userId, err := utils.ConvertStringToUUID(user.Data.Id)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	userWS := entity.NewUserWorkspace()
	userWS.Role = entity.ROLE_USER
	userWS.WorkspaceID = workspacePayload.WorkspaceID
	userWS.UserID = userId
	userWS.Email = data.Email
	userWS.PhoneNumber = data.PhoneNumber
	if err := s.postgresRepo.UserWorkspaceRepo.Create(ctx, tx, userWS); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.CreateUserWorkspaceResponse{}, nil
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
