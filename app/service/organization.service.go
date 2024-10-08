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
	payload, err := utils.GetGinContext[*utils.UserPayload](ctx, string(utils.USER_CONTEXT_KEY))
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
		WorkspaceId:          payload.Id,
		ParentOrganizationId: data.ParentOrganizationId,
		ParentLeaderId:       data.ParentOrganizationLeaderId,
		LeaderId:             data.LeaderId,
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
	organization, err := s.postgresRepo.OrganizationRepo.FindOneByFilter(ctx, nil, &repository.FindOrganizationByFilter{
		Id: &data.Id,
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
		LeaderId:             data.LeaderId,
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
	organization, err := s.postgresRepo.OrganizationRepo.FindOneByFilter(ctx, nil, &repository.FindOrganizationByFilter{
		Id: &data.Id,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeOrganizationNotFound)
	}

	// Check existed child organization
	existedChildOrganization, err := s.postgresRepo.OrganizationRepo.FindByFilter(ctx, nil, &repository.FindOrganizationByFilter{
		ParentOrganizationId: &data.Id,
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
		OrganizationId: &data.Id,
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
