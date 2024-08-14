package service

import (
	"context"
	"fmt"

	"workspace-server/app/entity"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	"workspace-server/app/repository"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/external/client"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	logger "workspace-server/package/log"
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
	userId := userPayload.Id.String()
	user, err := clientGRPC.GetUserByFilter(ctx, &oauth.GetUserByFilterParams{
		Id: &userId,
	})
	if err != nil {
		return nil, err
	}
	wsOwned, err := s.postgresRepo.UserWorkspaceRepo.CountByFilter(ctx, nil, &repository.FindUserWorkspaceByFilter{
		Role:     &role,
		IsActive: &isActive,
	})
	if err != nil {
		logger.Println(logger.LogPrintln{
			Ctx:       ctx,
			FileName:  "app/service/workspace.service.go",
			FuncName:  "CreateWorkspace",
			TraceData: "",
			Msg:       fmt.Sprintf("CountByFilter - %s", err.Error()),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if wsOwned >= int64(user.Data.GetLimitWorkspace()) {
		return nil, errors.New(errors.ErrCodePassedLimitWorkspace)
	}

	// Begin create workspace
	tx := database.BeginPostgresTransaction()
	ws := entity.NewWorkspace()
	ws.Email = data.Email
	ws.PhoneNumber = data.PhoneNumber
	ws.Address = data.Address
	ws.Name = data.Name

	// Create workspace
	if err := s.postgresRepo.WorkspaceRepo.Create(ctx, tx, ws); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	// Create user workspace
	userWS, err := s.helpers.UserWSHelper.CreateUserWS(ctx, &helper.CreateUserWsParams{
		Tx:          tx,
		UserId:      userPayload.Id,
		Email:       user.Data.Email,
		PhoneNumber: user.Data.PhoneNumber,
		Role:        entity.ROLE_OWNER,
		WorkspaceId: ws.Id,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// Create organization
	if err := s.helpers.OrganizationHelper.CreateOrganization(ctx, &helper.CreateOrganizationParams{
		Tx:                   tx,
		WorkspaceId:          ws.Id,
		ParentOrganizationId: nil,
		ParentLeaderId:       nil,
		LeaderId:             &userWS.Id,
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
	workspaceId, err := utils.ConvertStringToUUID(data.Id)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Check duplicate workspace name
	existedWS, err := s.postgresRepo.WorkspaceRepo.FindExistedByFilter(ctx, nil, &repository.FindWorkspaceByFilter{
		Id:   &workspaceId,
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
		Filter:     &repository.FindWorkspaceByFilter{Id: &workspaceId},
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeWorkspaceNotFound)
	}

	// Update workspace
	workspace.Email = data.Email
	workspace.PhoneNumber = data.PhoneNumber
	workspace.Name = data.Name
	if err := s.postgresRepo.WorkspaceRepo.Update(ctx, tx, workspace); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Update organization
	level := entity.ORGANiZATION_LEVEL_ROOT
	organization, err := s.postgresRepo.OrganizationRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter: &repository.FindOrganizationByFilter{
			WorkspaceId: &workspaceId,
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
	workspaceId, err := utils.ConvertStringToUUID(data.Id)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	tx := database.BeginPostgresTransaction()
	// Find workspace
	workspace, err := s.postgresRepo.WorkspaceRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter:     &repository.FindWorkspaceByFilter{Id: &workspaceId},
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
		WorkspaceId: &workspaceId,
	}); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.InactiveWorkspaceResponse{}, nil
}
