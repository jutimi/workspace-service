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
	generateParentOrganizationIds(parentOrganizationIds string, parentOrganizationId uuid.UUID) string
	createUserWorkspaceOrganization(ctx context.Context, data *CreateUserWorkspaceOrganizationParams) error
	validateParentOrganizationIds(parentOrganizationIds string) error
	validateLeaderIds(leaderIds string) error

	CreateOrganization(ctx context.Context, data *CreateOrganizationParams) error
	GetParentOrganizationIds(organizationIds string) []string
}
