package repository

import (
	"context"

	"workspace-server/app/entity"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, tx *gorm.DB, workspace *entity.Workspace) error
	Update(ctx context.Context, tx *gorm.DB, workspace *entity.Workspace) error
	Delete(ctx context.Context, tx *gorm.DB, workspace *entity.Workspace) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindWorkspaceByFilter) (*entity.Workspace, error)
	FindByFilter(ctx context.Context, tx *gorm.DB, filter *FindWorkspaceByFilter) ([]entity.Workspace, error)
	BulkCreate(ctx context.Context, tx *gorm.DB, workspaces []entity.Workspace) error
	CountByFilter(ctx context.Context, tx *gorm.DB, filter *FindWorkspaceByFilter) (int64, error)
	FindExistedByFilter(ctx context.Context, tx *gorm.DB, filter *FindWorkspaceByFilter) ([]entity.Workspace, error)
	FindOneByFilterForUpdate(ctx context.Context, filter *FindByFilterForUpdateParams) (*entity.Workspace, error)
}

type UserWorkspaceRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userWorkspace *entity.UserWorkspace) error
	Update(ctx context.Context, tx *gorm.DB, userWorkspace *entity.UserWorkspace) error
	Delete(ctx context.Context, tx *gorm.DB, userWorkspace *entity.UserWorkspace) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceByFilter) (*entity.UserWorkspace, error)
	FindByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceByFilter) ([]entity.UserWorkspace, error)
	BulkCreate(ctx context.Context, tx *gorm.DB, userWorkspaces []entity.UserWorkspace) error
	CountByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceByFilter) (int64, error)
	InActiveByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceByFilter) error
}

type UserWorkspaceDetailRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userWorkspaceDetail *entity.UserWorkspaceDetail) error
	Update(ctx context.Context, tx *gorm.DB, userWorkspaceDetail *entity.UserWorkspaceDetail) error
	Delete(ctx context.Context, tx *gorm.DB, userWorkspaceDetail *entity.UserWorkspaceDetail) error
	BulkCreate(ctx context.Context, tx *gorm.DB, userWorkspaceDetails []entity.UserWorkspaceDetail) error
}

type RedisRepository interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type OrganizationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, organization *entity.Organization) error
	Update(ctx context.Context, tx *gorm.DB, organization *entity.Organization) error
	Delete(ctx context.Context, tx *gorm.DB, organization *entity.Organization) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindOrganizationByFilter) (*entity.Organization, error)
	FindByFilter(ctx context.Context, tx *gorm.DB, filter *FindOrganizationByFilter) ([]entity.Organization, error)
	FindExistedByFilter(ctx context.Context, tx *gorm.DB, filter *FindOrganizationByFilter) ([]entity.Organization, error)
	FindByFilterForUpdate(ctx context.Context, data *FindByFilterForUpdateParams) ([]entity.Organization, error)
	FindOneByFilterForUpdate(ctx context.Context, data *FindByFilterForUpdateParams) (*entity.Organization, error)
}

type UserWorkspaceOrganizationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userWorkspaceOrganization *entity.UserWorkspaceOrganization) error
	Update(ctx context.Context, tx *gorm.DB, userWorkspaceOrganization *entity.UserWorkspaceOrganization) error
	Delete(ctx context.Context, tx *gorm.DB, userWorkspaceOrganization *entity.UserWorkspaceOrganization) error
	BulkCreate(ctx context.Context, tx *gorm.DB, userWorkspaceOrganizations []entity.UserWorkspaceOrganization) error
	FindByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceOrganizationFilter) ([]entity.UserWorkspaceOrganization, error)
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceOrganizationFilter) (*entity.UserWorkspaceOrganization, error)
	FindByFilterForUpdate(ctx context.Context, data *FindByFilterForUpdateParams) ([]entity.UserWorkspaceOrganization, error)
	DeleteByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceOrganizationFilter) error
}
