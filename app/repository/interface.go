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
	FindByFilter(ctx context.Context, tx *gorm.DB, filer *FindWorkspaceByFilter) ([]entity.Workspace, error)
	BulkCreate(ctx context.Context, tx *gorm.DB, workspaces []entity.Workspace) error
}

type UserWorkspaceRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userWorkspace *entity.UserWorkspace) error
	Update(ctx context.Context, tx *gorm.DB, userWorkspace *entity.UserWorkspace) error
	Delete(ctx context.Context, tx *gorm.DB, userWorkspace *entity.UserWorkspace) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserWorkspaceByFilter) (*entity.UserWorkspace, error)
	FindByFilter(ctx context.Context, tx *gorm.DB, filer *FindUserWorkspaceByFilter) ([]entity.UserWorkspace, error)
	BulkCreate(ctx context.Context, tx *gorm.DB, userWorkspaces []entity.UserWorkspace) error
}

type UserWorkspaceDetailRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userWorkspaceDetail *entity.UserWorkspaceDetail) error
	Update(ctx context.Context, tx *gorm.DB, userWorkspaceDetail *entity.UserWorkspaceDetail) error
	Delete(ctx context.Context, tx *gorm.DB, userWorkspaceDetail *entity.UserWorkspaceDetail) error
	BulkCreate(ctx context.Context, tx *gorm.DB, userWorkspaceDetails []entity.UserWorkspaceDetail) error
}
