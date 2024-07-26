package helper

import (
	"context"
	"workspace-server/app/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserWSHelper interface {
	CreateUserWS(ctx context.Context, data *CreateUserWsParams) (*entity.UserWorkspace, error)
}

type OrganizationHelper interface {
	generateParentOrganizationIds(parentOrganizationIds string, parentOrganizationId uuid.UUID) string
	createUserWorkspaceOrganization(ctx context.Context, data *CreateUserWorkspaceOrganizationParams) error
	validateParentOrganizationIds(ctx context.Context, tx *gorm.DB, parentOrganizationIds string) error
	validateLeaderIds(ctx context.Context, tx *gorm.DB, leaderIds string) error

	CreateOrganization(ctx context.Context, data *CreateOrganizationParams) error
	GetParentIds(parentIdStr string) []string
}
