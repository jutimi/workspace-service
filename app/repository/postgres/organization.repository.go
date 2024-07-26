package postgres_repository

import (
	"context"
	"errors"
	"time"
	"workspace-server/app/entity"
	"workspace-server/app/repository"
	"workspace-server/utils"

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
	organization.BaseWorkspace.UpdatedAt = time.Now().Unix()

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
	filter *repository.FindOrganizationByFilter,
) (*entity.Organization, error) {
	var data *entity.Organization
	query := r.buildFilter(ctx, nil, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *organizationRepository) FindByFilter(
	ctx context.Context,
	filter *repository.FindOrganizationByFilter,
) ([]entity.Organization, error) {
	var data []entity.Organization
	query := r.buildFilter(ctx, nil, filter)

	err := query.Find(&data).Error
	return data, err
}

func (r *organizationRepository) FindDuplicateOrganization(
	ctx context.Context,
	name string,
) (
	[]entity.Organization,
	error,
) {
	var data []entity.Organization
	nameSlug := utils.Slugify(name)

	query := r.db.WithContext(ctx).Where("name_slug = ?", nameSlug).Find(&data)
	if query.Error != nil {
		return nil, query.Error
	}
	return data, nil
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
		query = query.Scopes(findByString(*filter.ID, "email"))
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

	return query
}
