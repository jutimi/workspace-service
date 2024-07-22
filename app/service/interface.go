package service

import (
	"context"
	"workspace-server/app/model"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, data *model.CreateWorkspaceRequest) (*model.CreateWorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, data *model.UpdateWorkspaceRequest) (*model.UpdateWorkspaceResponse, error)
}
