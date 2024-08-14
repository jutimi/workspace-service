package postgres_repository

import (
	"context"
	"time"

	"workspace-server/app/entity"
	"workspace-server/app/repository"

	"gorm.io/gorm"
)

type userWorkspaceRepository struct {
	db *gorm.DB
}

func NewUserWorkspaceRepository(db *gorm.DB) repository.UserWorkspaceRepository {
	return &userWorkspaceRepository{
		db: db,
	}
}

func (r *userWorkspaceRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspace *entity.UserWorkspace,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&userWorkspace).Error
	}

	return r.db.WithContext(ctx).Create(&userWorkspace).Error
}

func (r *userWorkspaceRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspace *entity.UserWorkspace,
) error {
	userWorkspace.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&userWorkspace).Error
	}

	return r.db.WithContext(ctx).Save(&userWorkspace).Error
}

func (r *userWorkspaceRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspace *entity.UserWorkspace,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&userWorkspace).Error
	}

	return r.db.WithContext(ctx).Delete(&userWorkspace).Error
}

func (r *userWorkspaceRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	userWorkspaces []entity.UserWorkspace,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&userWorkspaces).Error
	}

	return r.db.WithContext(ctx).Create(&userWorkspaces).Error
}

func (r *userWorkspaceRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceByFilter,
) (*entity.UserWorkspace, error) {
	var data *entity.UserWorkspace
	query := r.buildFilter(ctx, tx, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *userWorkspaceRepository) FindByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceByFilter,
) ([]entity.UserWorkspace, error) {
	var data []entity.UserWorkspace
	query := r.buildFilter(ctx, tx, filter)

	err := query.Find(&data).Error
	return data, err
}

func (r *userWorkspaceRepository) CountByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceByFilter,
) (int64, error) {
	var count int64
	query := r.buildFilter(ctx, tx, filter)

	err := query.Model(&entity.UserWorkspace{}).Count(&count).Error
	return count, err
}

func (r *userWorkspaceRepository) InActiveByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceByFilter,
) error {
	query := r.buildFilter(ctx, nil, filter)
	return query.Model(entity.UserWorkspace{}).Updates(entity.UserWorkspace{IsActive: false}).Error
}

func (r *userWorkspaceRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Scopes(findByText(*filter.Email, "email"))
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Scopes(findByText(*filter.PhoneNumber, "phone_number"))
	}
	if filter.Id != nil {
		query = query.Scopes(findByString(*filter.Id, "id"))
	}
	if filter.Ids != nil && len(filter.Ids) > 0 {
		query = query.Scopes(findBySlice(filter.Ids, "id"))
	}
	if filter.Emails != nil && len(filter.Emails) > 0 {
		query = query.Scopes(findBySlice(filter.Emails, "email"))
	}
	if filter.PhoneNumbers != nil && len(filter.PhoneNumbers) > 0 {
		query = query.Scopes(findBySlice(filter.PhoneNumbers, "phone_number"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.Role != nil {
		query = query.Scopes(findByString(*filter.Role, "role"))
	}

	// Relation query
	if filter.IsIncludeDetail {
		query = query.Preload("UserWorkspaceDetail")
	} else if filter.IsRequireDetail {
		query = query.InnerJoins("UserWorkspaceUser", func(db *gorm.DB) *gorm.DB {
			return r.buildRelationFilter(db, filter)
		})
	}

	return query
}

func (r *userWorkspaceRepository) buildRelationFilter(
	db *gorm.DB,
	filter *repository.FindUserWorkspaceByFilter,
) *gorm.DB {
	query := db

	//
	if filter.Name != nil {
		query.Scopes(findByName(*filter.Name, "full_name_slug"))
	}

	return query
}
