package postgres_repository

import (
	"context"
	"time"

	"workspace-server/app/entity"
	"workspace-server/app/repository"

	"gorm.io/gorm"
)

type userWorkspaceDetailRepository struct {
	db *gorm.DB
}

func NewUserWorkspaceDetailRepository(db *gorm.DB) repository.UserWorkspaceDetailRepository {
	return &userWorkspaceDetailRepository{
		db: db,
	}
}

func (r *userWorkspaceDetailRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceDetail *entity.UserWorkspaceDetail,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&userWorkspaceDetail).Error
	}

	return r.db.WithContext(ctx).Create(&userWorkspaceDetail).Error
}

func (r *userWorkspaceDetailRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceDetail *entity.UserWorkspaceDetail,
) error {
	userWorkspaceDetail.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&userWorkspaceDetail).Error
	}

	return r.db.WithContext(ctx).Save(&userWorkspaceDetail).Error
}

func (r *userWorkspaceDetailRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceDetail *entity.UserWorkspaceDetail,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&userWorkspaceDetail).Error
	}

	return r.db.WithContext(ctx).Delete(&userWorkspaceDetail).Error
}

func (r *userWorkspaceDetailRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaceDetails []entity.UserWorkspaceDetail,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&userWorkspaceDetails).Error
	}

	return r.db.WithContext(ctx).Create(&userWorkspaceDetails).Error
}
