package service

import (
	"context"

	"github.com/jutimi/workspace-server/app/helper"
	"github.com/jutimi/workspace-server/app/model"
	"github.com/jutimi/workspace-server/app/repository"
	postgres_repository "github.com/jutimi/workspace-server/app/repository/postgres"
	"github.com/jutimi/workspace-server/package/database"
	"github.com/jutimi/workspace-server/package/errors"
	"github.com/jutimi/workspace-server/utils"
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
	// Check organization
	organization, err := s.postgresRepo.OrganizationRepo.FindOneByFilter(ctx, &repository.FindOrganizationByFilter{
		ID: &data.Id,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeOrganizationNotFound)
	}

	// Begin update organization
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := s.helpers.OrganizationHelper.UpdateOrganization(ctx, &helper.UpdateOrganizationParams{
		Tx:                   tx,
		Organization:         organization,
		ParentOrganizationId: data.ParentOrganizationId,
		ParentLeaderId:       data.ParentOrganizationLeaderId,
		LeaderID:             data.LeaderId,
		Name:                 data.Name,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.UpdateOrganizationResponse{}, nil
}

func (s *organizationService) RemoveOrganization(
	ctx context.Context,
	data *model.RemoveOrganizationRequest,
) (*model.RemoveOrganizationResponse, error) {
	limit := 1
	offset := 0

	// Check organization
	organization, err := s.postgresRepo.OrganizationRepo.FindOneByFilter(ctx, &repository.FindOrganizationByFilter{
		ID: &data.Id,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeOrganizationNotFound)
	}

	// Check existed child organization
	existedChildOrganization, err := s.postgresRepo.OrganizationRepo.FindByFilter(ctx, &repository.FindOrganizationByFilter{
		ParentOrganizationID: &data.Id,
		Limit:                &limit,
		Offset:               &offset,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedChildOrganization) > 0 {
		return nil, errors.New(errors.ErrCodeOrganizationHasChild)
	}

	// Begin organization
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Remove user workspace organization
	if err := s.postgresRepo.UserWorkspaceOrganizationRepo.DeleteByFilter(ctx, tx, &repository.FindUserWorkspaceOrganizationFilter{
		OrganizationID: &data.Id,
	}); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Remove organization
	if err := s.postgresRepo.OrganizationRepo.Delete(ctx, tx, organization); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.RemoveOrganizationResponse{}, nil
}
