package service

import (
	"context"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/package/database"
	"workspace-server/package/errors"
	"workspace-server/utils"
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

	// Create organization
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := s.helpers.OrganizationHelper.CreateOrganization(ctx, &helper.CreateOrganizationParams{
		Tx:                   tx,
		WorkspaceID:          payload.ID,
		ParentOrganizationId: data.ParentOrganizationId,
		ParentLeaderId:       data.ParentOrganizationLeaderId,
		LeaderID:             data.LeaderId,
		Name:                 data.Name,
		SubLeaders:           nil,
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
