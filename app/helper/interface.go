package helper

import (
	"context"

	"workspace-server/app/entity"
	"workspace-server/app/model"

	"github.com/google/uuid"
)

type UserWorkspaceHelper interface {
	CreateUserWorkspace(ctx context.Context, data *CreateUserWorkspaceParams) (*entity.UserWorkspace, error)
}

type OrganizationHelper interface {
	generateParentIds(ctx context.Context, parentOrganizationIds string, parentOrganizationId uuid.UUID) string
	createUserWorkspaceOrganization(ctx context.Context, data *CreateUserWorkspaceOrganizationParams) error
	validateParentOrganizationIds(ctx context.Context, parentOrganizationIds string) error
	validateParentOrganization(ctx context.Context, data *validateOrganizationParams) (*validateOrganizationResult, error)
	validateUpdateOrganization(ctx context.Context, data *UpdateOrganizationParams) error
	validateDuplicateUserWorkspace(ctx context.Context, leaderId *uuid.UUID, subLeaders []model.SubLeaderData) error

	CreateOrganization(ctx context.Context, data *CreateOrganizationParams) error
	UpdateOrganization(ctx context.Context, data *UpdateOrganizationParams) error
	GetParentIds(parentIdStr string) []uuid.UUID
}
