package postgres_repository

import (
	"context"
	"errors"
	"time"

	"workspace-server/app/entity"
	"workspace-server/app/repository"

	"gorm.io/gorm"
)

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) repository.WorkspaceRepository {
	return &workspaceRepository{
		db: db,
	}
}

func (r *workspaceRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	workspace *entity.Workspace,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&workspace).Error
	}

	return r.db.WithContext(ctx).Create(&workspace).Error
}

func (r *workspaceRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	workspace *entity.Workspace,
) error {
	workspace.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&workspace).Error
	}

	return r.db.WithContext(ctx).Save(&workspace).Error
}

func (r *workspaceRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	workspace *entity.Workspace,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&workspace).Error
	}

	return r.db.WithContext(ctx).Delete(&workspace).Error
}

func (r *workspaceRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	workspaces []entity.Workspace,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&workspaces).Error
	}

	return r.db.WithContext(ctx).Create(&workspaces).Error
}

func (r *workspaceRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWorkspaceByFilter,
) (*entity.Workspace, error) {
	var data *entity.Workspace
	query := r.buildFilter(ctx, tx, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *workspaceRepository) FindByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWorkspaceByFilter,
) ([]entity.Workspace, error) {
	var data []entity.Workspace
	query := r.buildFilter(ctx, tx, filter)

	err := query.Find(&data).Error
	return data, err
}

func (r *workspaceRepository) CountByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWorkspaceByFilter,
) (int64, error) {
	var count int64
	query := r.buildFilter(ctx, tx, filter)

	err := query.Count(&count).Error
	return count, err
}

func (r *workspaceRepository) FindExistedByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWorkspaceByFilter,
) ([]entity.Workspace, error) {
	var data []entity.Workspace
	query := r.buildExistedFilter(ctx, tx, filter)

	err := query.Find(&data).Error
	return data, err
}

func (r *workspaceRepository) FindOneByFilterForUpdate(
	ctx context.Context,
	data *repository.FindByFilterForUpdateParams,
) (*entity.Workspace, error) {
	filter, ok := data.Filter.(*repository.FindWorkspaceByFilter)
	if !ok {
		return nil, errors.New("invalid argument")
	}

	var workspace *entity.Workspace
	query := r.buildFilter(ctx, data.Tx, filter)
	query = buildLockQuery(query, data.LockOption)
	err := query.First(&workspace).Error
	return workspace, err
}

// ------------------------------------------------------------------------------
func (r *workspaceRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWorkspaceByFilter,
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
	if filter.ID != nil {
		query = query.Scopes(findByString(*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(findBySlice(filter.IDs, "id"))
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
	if filter.Name != nil {
		query = query.Scopes(findByName(*filter.Name, "name_slug"))
	}

	return query
}

func (r *workspaceRepository) buildExistedFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWorkspaceByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Scopes(orByText(*filter.Email, "email"))
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Scopes(orByText(*filter.PhoneNumber, "phone_number"))
	}
	if filter.ID != nil {
		query = query.Scopes(excludeByString(*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(excludeBySlice(filter.IDs, "id"))
	}
	if filter.Emails != nil && len(filter.Emails) > 0 {
		query = query.Scopes(orBySlice(filter.Emails, "email"))
	}
	if filter.PhoneNumbers != nil && len(filter.PhoneNumbers) > 0 {
		query = query.Scopes(orBySlice(filter.PhoneNumbers, "phone_number"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.Name != nil {
		query = query.Scopes(orByName(*filter.Name, "name_slug"))
	}

	return query
}
