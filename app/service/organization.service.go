package service

import (
	"context"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	"workspace-server/app/repository"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/google/uuid"
)

type organizationService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewOrganizationService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) OrganizationService {
	return &organizationService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
	}
}

func (s *organizationService) CreateOrganization(
	ctx context.Context,
	data *model.CreateOrganizationRequest,
) (*model.CreateOrganizationResponse, error) {
	payload, err := utils.GetScopeContext[*utils.UserPayload](ctx, string(utils.USER_CONTEXT_KEY))
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Get user workspace data
	userWSIds := make([]uuid.UUID, 0)
	if data.LeaderId != nil && *data.LeaderId != "" {
		leaderId, err := utils.ConvertStringToUUID(*data.LeaderId)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		userWSIds = append(userWSIds, leaderId)
	}
	for _, leaderId := range data.SubLeaderIds {
		id, err := utils.ConvertStringToUUID(leaderId)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		userWSIds = append(userWSIds, id)
	}
	userWS, err := s.postgresRepo.UserWorkspaceRepo.FindByFilter(ctx, &repository.FindUserWorkspaceByFilter{
		IDs: userWSIds,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(userWS) != len(userWSIds) {
		return nil, errors.New(errors.ErrCodeUserWorkspaceNotFound)
	}

	// Create organization
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := s.helpers.OrganizationHelper.CreateOrganization(ctx, &helper.CreateOrganizationParams{
		Tx:                 tx,
		WorkspaceID:        payload.ID,
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

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.CreateOrganizationResponse{}, nil
}

func (s *organizationService) UpdateOrganization(
	ctx context.Context,
	data *model.UpdateOrganizationRequest,
) (*model.UpdateOrganizationResponse, error) {
	return nil, nil
}

func (s *organizationService) RemoveOrganization(
	ctx context.Context,
	data *model.RemoveOrganizationRequest,
) (*model.RemoveOrganizationResponse, error) {
	return nil, nil
}
