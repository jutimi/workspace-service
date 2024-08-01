package helper

import (
	"context"
	"workspace-server/app/entity"

	"github.com/google/uuid"
)

type UserWSHelper interface {
	CreateUserWS(ctx context.Context, data *CreateUserWsParams) (*entity.UserWorkspace, error)
}

type OrganizationHelper interface {
	generateParentIds(parentOrganizationIds string, parentOrganizationId uuid.UUID) string
	createUserWorkspaceOrganization(ctx context.Context, data *CreateUserWorkspaceOrganizationParams) error
	validateParentOrganizationIds(ctx context.Context, parentOrganizationIds string) error
	validateOrganization(ctx context.Context, data *validateOrganizationParams) (*validateOrganizationResult, error)

	CreateOrganization(ctx context.Context, data *CreateOrganizationParams) error
	UpdateOrganization(ctx context.Context, data *UpdateOrganizationParams) error
	GetParentIds(parentIdStr string) []uuid.UUID
}
