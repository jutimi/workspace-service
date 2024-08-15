package postgres_repository

import (
	"context"
	"errors"
	"time"

	"workspace-server/app/entity"
	"workspace-server/app/repository"

	"github.com/google/uuid"
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
	userWorkspaceOrganization.UpdatedAt = time.Now().Unix()

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
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceOrganizationFilter,
) ([]entity.UserWorkspaceOrganization, error) {
	var userWorkspaceOrganizations []entity.UserWorkspaceOrganization

	err := r.buildFilter(ctx, tx, filter).Find(&userWorkspaceOrganizations).Error
	return userWorkspaceOrganizations, err
}

func (r *userWorkspaceOrganizationRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserWorkspaceOrganizationFilter,
) (*entity.UserWorkspaceOrganization, error) {
	var userWorkspaceOrganizations *entity.UserWorkspaceOrganization

	err := r.buildFilter(ctx, tx, filter).First(&userWorkspaceOrganizations).Error
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

	if filter.Id != nil && *filter.Id != uuid.Nil {
		query = query.Scopes(whereBy(*filter.Id, "id"))
	}
	if filter.Ids != nil && len(filter.Ids) > 0 {
		query = query.Scopes(whereBySlice(filter.Ids, "id"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.LeaderIds != nil && len(*filter.LeaderIds) > 0 {
		query = query.Scopes(whereByText(*filter.LeaderIds, "leader_ids"))
	}

	// Relation query
	if filter.IsIncludeOrganization {
		query = query.Preload("Organization")
	} else if filter.IsRequireOrganization {
		query = query.InnerJoins("Organization")
	}

	return query
}
