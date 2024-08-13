package service

import (
	"context"

	"workspace-server/app/entity"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	"workspace-server/app/repository"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/external/client"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/jutimi/grpc-service/oauth"
	"gorm.io/gorm/clause"
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
	clientGRPC := client.NewOAuthClient()
	defer clientGRPC.CloseConn()

	// Get user payload from context
	userPayload, err := utils.GetGinContext[*utils.UserPayload](ctx, string(utils.USER_CONTEXT_KEY))
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Get user data and check limit workspace
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

	// Get email and phone number
	email := user.Data.Email
	phoneNumber := user.Data.PhoneNumber
	if data.Email != nil && *data.Email != "" {
		email = data.Email
	}
	if data.PhoneNumber != nil && *data.PhoneNumber != "" {
		phoneNumber = data.PhoneNumber
	}

	// Begin create workspace
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
		Tx:                   tx,
		WorkspaceID:          ws.ID,
		ParentOrganizationId: nil,
		ParentLeaderId:       nil,
		LeaderID:             &userWS.ID,
		Name:                 data.Name,
		SubLeaders:           nil,
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
	workspaceId, err := utils.ConvertStringToUUID(data.ID)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Check duplicate workspace name
	existedWS, err := s.postgresRepo.WorkspaceRepo.FindExistedByFilter(ctx, &repository.FindWorkspaceByFilter{
		ID:   &workspaceId,
		Name: &data.Name,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedWS) > 0 {
		return nil, errors.New(errors.ErrCodeDuplicateWorkspaceName)
	}

	tx := database.BeginPostgresTransaction()
	// Find workspace
	workspace, err := s.postgresRepo.WorkspaceRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter:     &repository.FindWorkspaceByFilter{ID: &workspaceId},
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeWorkspaceNotFound)
	}

	// Update workspace
	if data.Email != nil && *data.Email != "" {
		workspace.Email = *data.Email
	}
	if data.PhoneNumber != nil && *data.PhoneNumber != "" {
		workspace.PhoneNumber = *data.PhoneNumber
	}
	if err := s.postgresRepo.WorkspaceRepo.Update(ctx, tx, workspace); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Update organization
	level := entity.ORGANiZATION_LEVEL_ROOT
	organization, err := s.postgresRepo.OrganizationRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter: &repository.FindOrganizationByFilter{
			WorkspaceID: &workspaceId,
			Level:       &level,
		},
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	organization.Name = workspace.Name
	if err := s.postgresRepo.OrganizationRepo.Update(ctx, tx, organization); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	tx.Commit()

	return &model.UpdateWorkspaceResponse{}, nil
}

func (s *workspaceService) InactiveWorkspace(
	ctx context.Context,
	data *model.InactiveWorkspaceRequest,
) (*model.InactiveWorkspaceResponse, error) {
	workspaceId, err := utils.ConvertStringToUUID(data.ID)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	tx := database.BeginPostgresTransaction()
	// Find workspace
	workspace, err := s.postgresRepo.WorkspaceRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter:     &repository.FindWorkspaceByFilter{ID: &workspaceId},
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeWorkspaceNotFound)
	}

	// Update workspace inactive
	workspace.IsActive = false
	if err := s.postgresRepo.WorkspaceRepo.Update(ctx, tx, workspace); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	tx.Commit()

	// Update user workspace inactive
	if err := s.postgresRepo.UserWorkspaceRepo.InActiveByFilter(ctx, tx, &repository.FindUserWorkspaceByFilter{
		WorkspaceID: &workspaceId,
	}); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.InactiveWorkspaceResponse{}, nil
}
