package postgres_repository

import (
	"context"
	"errors"
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

func (r *userWorkspaceOrganizationRepository) FindByFilter(
	ctx context.Context,
	filter *repository.FindUserWorkspaceOrganizationFilter,
) ([]entity.UserWorkspaceOrganization, error) {
	var userWorkspaceOrganizations []entity.UserWorkspaceOrganization

	err := r.buildFilter(ctx, nil, filter).Find(&userWorkspaceOrganizations).Error
	return userWorkspaceOrganizations, err
}

func (r *userWorkspaceOrganizationRepository) FindOneByFilter(
	ctx context.Context,
	filter *repository.FindUserWorkspaceOrganizationFilter,
) (*entity.UserWorkspaceOrganization, error) {
	var userWorkspaceOrganizations *entity.UserWorkspaceOrganization

	err := r.buildFilter(ctx, nil, filter).First(&userWorkspaceOrganizations).Error
	return userWorkspaceOrganizations, err
}

func (r *userWorkspaceOrganizationRepository) FindByFilterForUpdate(
	ctx context.Context,
	data *repository.FindByFilterForUpdateParams,
) ([]entity.UserWorkspaceOrganization, error) {
	filter, ok := data.Filter.(*repository.FindUserWorkspaceOrganizationFilter)
	if !ok {
		return nil, errors.New("invalid argument")
	}

	var users []entity.UserWorkspaceOrganization
	query := r.buildFilter(ctx, data.Tx, filter)
	query = buildLockQuery(query, data.LockOption)

	err := query.Find(&users).Error
	return users, err
}

func (r *userWorkspaceOrganizationRepository) DeleteByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceOrganizationFilter,
) error {
	query := r.buildFilter(ctx, tx, filter)
	return query.Delete(&entity.UserWorkspaceOrganization{}).Error
}

// ------------------------------------------------------------------------------
func (r *userWorkspaceOrganizationRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceOrganizationFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.ID != nil {
		query = query.Scopes(findByString(*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(findBySlice(filter.IDs, "id"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.LeaderIDs != nil && len(*filter.LeaderIDs) > 0 {
		query = query.Scopes(findByText(*filter.LeaderIDs, "leader_ids"))
	}

	// Relation query
	if filter.IsIncludeOrganization {
		query = query.Preload("Organization")
	} else if filter.IsRequireOrganization {
		query = query.InnerJoins("Organization")
	}

	return query
}
