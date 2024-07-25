package service

import (
	"context"
	"workspace-server/app/entity"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	"workspace-server/app/repository"
	postgres_repository "workspace-server/app/repository/postgres"
	client_grpc "workspace-server/grpc/client"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/jutimi/grpc-service/oauth"
)

type workspaceService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewWorkspaceService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) WorkspaceService {
	return &workspaceService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
	}
}

func (s *workspaceService) CreateWorkspace(
	ctx context.Context,
	data *model.CreateWorkspaceRequest,
) (*model.CreateWorkspaceResponse, error) {
	isActive := true
	role := entity.ROLE_OWNER
	clientGRPC := client_grpc.NewOAuthClient()
	defer clientGRPC.CloseConn()

	userPayload, err := utils.GetScopeContext[*utils.UserPayload](ctx, string(utils.USER_CONTEXT_KEY))
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	userId := userPayload.ID.String()
	user, err := clientGRPC.GetUserByFilter(ctx, &oauth.GetUserByFilterParams{
		Id: &userId,
	})
	if err != nil {
		return nil, err
	}
	wsOwned, err := s.postgresRepo.UserWorkspaceRepo.CountByFilter(ctx, &repository.FindUserWorkspaceByFilter{
		Role:     &role,
		IsActive: &isActive,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if wsOwned >= int64(user.Data.GetLimitWorkspace()) {
		return nil, errors.New(errors.ErrCodePassedLimitWorkspace)
	}

	email := user.Data.Email
	phoneNumber := user.Data.PhoneNumber
	if data.Email != nil && *data.Email != "" {
		email = data.Email
	}
	if data.PhoneNumber != nil && *data.PhoneNumber != "" {
		phoneNumber = data.PhoneNumber
	}

	tx := database.BeginPostgresTransaction()
	ws := entity.NewWorkspace()
	ws.Email = *email
	ws.PhoneNumber = *phoneNumber
	ws.Address = data.Address

	// Create workspace
	if err := s.postgresRepo.WorkspaceRepo.Create(ctx, tx, ws); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	// Create user workspace
	userWS, err := s.helpers.UserWSHelper.CreateUserWS(ctx, &helper.CreateUserWsParams{
		Tx:          tx,
		UserID:      userPayload.ID,
		Email:       user.Data.Email,
		PhoneNumber: user.Data.PhoneNumber,
		Role:        entity.ROLE_OWNER,
		WorkspaceID: ws.ID,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create organization
	if err := s.helpers.OrganizationHelper.CreateOrganization(ctx, &helper.CreateOrganizationParams{
		Tx:                 tx,
		WorkspaceID:        ws.ID,
		ParentOrganization: nil,
		Leader: &helper.MemberInfo{
			LeaderIds:     nil,
			UserWorkspace: *userWS,
		},
		Name:       data.Name,
		SubLeaders: nil,
		Members:    nil,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &model.CreateWorkspaceResponse{}, nil
}

func (s *workspaceService) UpdateWorkspace(
	ctx context.Context,
	data *model.UpdateWorkspaceRequest,
) (*model.UpdateWorkspaceResponse, error) {
	return nil, nil
}

func (s *workspaceService) InactiveWorkspace(
	ctx context.Context,
	data *model.InactiveWorkspaceRequest,
) (*model.InactiveWorkspaceResponse, error) {
	return nil, nil
}
