package service

import (
	"context"
	"workspace-server/app/helper"
	"workspace-server/app/model"
	postgres_repository "workspace-server/app/repository/postgres"
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
	return nil, nil
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
