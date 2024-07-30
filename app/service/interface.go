package service

import (
	"context"
	"workspace-server/app/model"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, data *model.CreateWorkspaceRequest) (*model.CreateWorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, data *model.UpdateWorkspaceRequest) (*model.UpdateWorkspaceResponse, error)
	InactiveWorkspace(ctx context.Context, data *model.InactiveWorkspaceRequest) (*model.InactiveWorkspaceResponse, error)
}

type UserWorkspaceService interface {
	CreateUserWorkspace(ctx context.Context, data *model.CreateUserWorkspaceRequest) (*model.CreateUserWorkspaceResponse, error)
	UpdateUserWorkspace(ctx context.Context, data *model.UpdateUserWorkspaceRequest) (*model.UpdateUserWorkspaceResponse, error)
	InactiveUserWorkspace(ctx context.Context, data *model.InactiveUserWorkspaceRequest) (*model.InactiveUserWorkspaceResponse, error)
	RemoveUserWorkspace(ctx context.Context, data *model.RemoveUserWorkspaceRequest) (*model.RemoveUserWorkspaceResponse, error)
}

type OrganizationService interface {
	CreateOrganization(ctx context.Context, data *model.CreateOrganizationRequest) (*model.CreateOrganizationResponse, error)
	UpdateOrganization(ctx context.Context, data *model.UpdateOrganizationRequest) (*model.UpdateOrganizationResponse, error)
	RemoveOrganization(ctx context.Context, data *model.RemoveOrganizationRequest) (*model.RemoveOrganizationResponse, error)
}
