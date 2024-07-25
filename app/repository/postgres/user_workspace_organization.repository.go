package postgres_repository

import (
	"context"
	"time"
	"workspace-server/app/entity"
	"workspace-server/app/repository"

	"gorm.io/gorm"
)

type userWorkspaceOrganizationRepository struct {
	db *gorm.DB
}

func NewUserWorkspaceOrganizationRepository(db *gorm.DB) repository.UserWorkspaceOrganizationRepository {
	return &userWorkspaceOrganizationRepository{
		db: db,
	}
}

func (r *userWorkspaceOrganizationRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceOrganization *entity.UserWorkspaceOrganization,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&userWorkspaceOrganization).Error
	}

	return r.db.WithContext(ctx).Create(&userWorkspaceOrganization).Error
}

func (r *userWorkspaceOrganizationRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceOrganization *entity.UserWorkspaceOrganization,
) error {
	userWorkspaceOrganization.BaseUserWorkspace.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&userWorkspaceOrganization).Error
	}

	return r.db.WithContext(ctx).Save(&userWorkspaceOrganization).Error
}

func (r *userWorkspaceOrganizationRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceOrganization *entity.UserWorkspaceOrganization,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&userWorkspaceOrganization).Error
	}

	return r.db.WithContext(ctx).Delete(&userWorkspaceOrganization).Error
}

func (r *userWorkspaceOrganizationRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceOrganizations []entity.UserWorkspaceOrganization,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&userWorkspaceOrganizations).Error
	}

	return r.db.WithContext(ctx).Create(&userWorkspaceOrganizations).Error
}
