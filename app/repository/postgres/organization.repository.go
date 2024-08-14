package postgres_repository

import (
	"context"
	"errors"
	"time"

	"workspace-server/app/entity"
	"workspace-server/app/repository"

	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

func (r *organizationRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	organization *entity.Organization,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&organization).Error
	}

	return r.db.WithContext(ctx).Create(&organization).Error
}

func (r *organizationRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	organization *entity.Organization,
) error {
	organization.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&organization).Error
	}

	return r.db.WithContext(ctx).Save(&organization).Error
}

func (r *organizationRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	organization *entity.Organization,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&organization).Error
	}

	return r.db.WithContext(ctx).Delete(&organization).Error
}

func (r *organizationRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOrganizationByFilter,
) (*entity.Organization, error) {
	var data *entity.Organization
	query := r.buildFilter(ctx, tx, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *organizationRepository) FindByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOrganizationByFilter,
) ([]entity.Organization, error) {
	var data []entity.Organization
	query := r.buildFilter(ctx, tx, filter)

	err := query.Find(&data).Error
	return data, err
}

func (r *organizationRepository) FindExistedByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOrganizationByFilter,
) (
	[]entity.Organization,
	error,
) {
	var data []entity.Organization
	query := r.buildExistedFilter(ctx, tx, filter)

	err := query.Find(&data).Error
	return data, err
}

func (r *organizationRepository) FindOneByFilterForUpdate(
	ctx context.Context,
	data *repository.FindByFilterForUpdateParams,
) (*entity.Organization, error) {
	filter, ok := data.Filter.(*repository.FindOrganizationByFilter)
	if !ok {
		return nil, errors.New("invalid argument")
	}

	var organization *entity.Organization
	query := r.buildFilter(ctx, data.Tx, filter)
	query = buildLockQuery(query, data.LockOption)

	err := query.First(&organization).Error
	return organization, err
}

func (r *organizationRepository) FindByFilterForUpdate(
	ctx context.Context,
	data *repository.FindByFilterForUpdateParams,
) ([]entity.Organization, error) {
	filter, ok := data.Filter.(*repository.FindOrganizationByFilter)
	if !ok {
		return nil, errors.New("invalid argument")
	}

	var organizations []entity.Organization
	query := r.buildFilter(ctx, data.Tx, filter)
	query = buildLockQuery(query, data.LockOption)

	err := query.Find(&organizations).Error
	return organizations, err
}

// ------------------------------------------------------------------------------
func (r *organizationRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOrganizationByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.ID != nil {
		query = query.Scopes(findByString(*filter.ID, "id"))
	}
	if filter.ParentOrganizationID != nil {
		query = query.Scopes(findByString(*filter.ParentOrganizationID, "parent_organization_id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(findBySlice(filter.IDs, "id"))
	}
	if filter.Name != nil {
		query = query.Scopes(findByName(*filter.Name, "name_slug"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.WorkspaceID != nil {
		query = query.Scopes(findByString(*filter.WorkspaceID, "workspace_id"))
	}
	if filter.Level != nil {
		query = query.Scopes(findByString(*filter.Level, "level"))
	}
	if filter.ParentOrganizationIDsStr != nil && len(*filter.ParentOrganizationIDsStr) > 0 {
		query = query.Scopes(findByText(*filter.ParentOrganizationIDsStr, "parent_organization_ids"))
	}

	return query
}

func (r *organizationRepository) buildExistedFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOrganizationByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.ID != nil {
		query = query.Scopes(excludeByString(*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(excludeBySlice(filter.IDs, "id"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.Name != nil {
		query = query.Scopes(orByName(*filter.Name, "name_slug"))
	}

	return query
}
