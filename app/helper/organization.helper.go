package helper

import (
	"context"
	"fmt"
	"strings"
	"workspace-server/app/entity"
	"workspace-server/app/repository"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type organizationHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewOrganizationHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) OrganizationHelper {
	return &organizationHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *organizationHelper) CreateOrganization(ctx context.Context, data *CreateOrganizationParams) error {
	leaderIds := make([]uuid.UUID, 0)

	// Create organization
	organization := entity.NewOrganization()
	organization.Name = data.Name
	organization.BaseWorkspace.WorkspaceID = data.WorkspaceID
	organization.Level = entity.ORGANiZATION_LEVEL_ROOT
	// Check parent organization
	if data.ParentOrganization != nil {
		organization.Level = data.ParentOrganization.Level + 1
		organization.ParentOrganizationID = &data.ParentOrganization.ID
		parentOrganizationIds := h.generateParentOrganizationIds(
			*data.ParentOrganization.ParentOrganizationIDs,
			data.ParentOrganization.ID,
		)
		organization.ParentOrganizationIDs = &parentOrganizationIds
	}
	if err := h.postgresRepo.OrganizationRepo.Create(ctx, data.Tx, organization); err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	// Create user workspace organization
	if data.Leader != nil {
		if err := h.createUserWorkspaceOrganization(ctx, &CreateUserWorkspaceOrganizationParams{
			Organization: organization,
			Role:         entity.ORGANiZATION_ROLE_LEADER,
			Data:         []MemberInfo{*data.Leader},
			Tx:           data.Tx,
		}); err != nil {
			return errors.New(errors.ErrCodeInternalServerError)
		}
	}
	if data.Members != nil {
		if err := h.createUserWorkspaceOrganization(ctx, &CreateUserWorkspaceOrganizationParams{
			Organization: organization,
			Role:         entity.ORGANiZATION_ROLE_MEMBER,
			Data:         data.Members,
			Tx:           data.Tx,
		}); err != nil {
			return errors.New(errors.ErrCodeInternalServerError)
		}
	}
	if data.SubLeaders != nil {
		if err := h.createUserWorkspaceOrganization(ctx, &CreateUserWorkspaceOrganizationParams{
			Organization: organization,
			Role:         entity.ORGANiZATION_ROLE_SUB_LEADER,
			Data:         data.SubLeaders,
			Tx:           data.Tx,
		}); err != nil {
			return errors.New(errors.ErrCodeInternalServerError)
		}
	}

	return nil
}

func (h *organizationHelper) GetParentIds(parentIdStr string) []string {
	return strings.Split(parentIdStr, "/")
}

// ----------------------------------------------- Helper -----------------------------------------------

func (h *organizationHelper) generateParentOrganizationIds(
	parentOrganizationIds string,
	parentOrganizationId uuid.UUID,
) string {
	return fmt.Sprintf("%s/%s", parentOrganizationIds, parentOrganizationId.String())
}

func (h *organizationHelper) createUserWorkspaceOrganization(
	ctx context.Context,
	data *CreateUserWorkspaceOrganizationParams,
) error {
	if data.Organization == nil {
		return nil
	}

	users := make([]entity.UserWorkspaceOrganization, 0)
	for _, userWSId := range data.UserWorkspaceIds {
		// Validate leader tree before create
		if err := h.validateLeaderIds(ctx, data.Tx, *member.LeaderIds); err != nil {
			return err
		}

		userWorkspaceOrganization := entity.NewUserWorkspaceOrganization()
		userWorkspaceOrganization.OrganizationID = data.Organization.ID
		userWorkspaceOrganization.BaseUserWorkspace.WorkspaceID = member.WorkspaceID
		userWorkspaceOrganization.BaseUserWorkspace.UserWorkspaceID = member.ID
		userWorkspaceOrganization.Role = data.Role
		userWorkspaceOrganization.LeaderIDs = member.LeaderIds
		users = append(users, *userWorkspaceOrganization)
	}
	if err := h.postgresRepo.UserWorkspaceOrganizationRepo.BulkCreate(ctx, data.Tx, users); err != nil {
		return err
	}

	return nil
}

func (h *organizationHelper) validateParentOrganizationIds(
	ctx context.Context,
	tx *gorm.DB,
	parentOrganizationIds string,
) error {
	var organizations []entity.Organization

	//
	organizationIds := h.GetParentIds(parentOrganizationIds)
	if len(organizationIds) == 0 {
		return nil
	}
	convertOrganizationIds, err := utils.ConvertSliceStringToUUID(organizationIds)
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	// Get organizations
	if tx != nil {
		organizations, err = h.postgresRepo.OrganizationRepo.FindByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
			Filter: &repository.FindOrganizationByFilter{
				IDs: convertOrganizationIds,
			},
			LockOption: clause.LockingStrengthShare,
			Tx:         tx,
		})
	} else {
		organizations, err = h.postgresRepo.OrganizationRepo.FindByFilter(ctx, &repository.FindOrganizationByFilter{
			IDs: convertOrganizationIds,
		})
	}
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(organizations) != len(organizationIds) {
		return errors.New(errors.ErrCodeOrganizationNotFound)
	}

	/*
		Check format of parentOrganizationIds. Rule:
		- ParentOrganizationIds level must be in ASC order
	*/
	level := entity.ORGANiZATION_LEVEL_ROOT - 1
	for _, organizationId := range convertOrganizationIds {
		for _, organization := range organizations {
			if organizationId != organization.ID {
				continue
			}
			if organization.Level <= level {
				return errors.New(errors.ErrCodeInvalidParentOrganizationIds)
			}

			level = organization.Level
		}
	}

	return nil
}

func (h *organizationHelper) validateLeaderIds(
	ctx context.Context,
	tx *gorm.DB,
	leaderIds string,
) error {
	var userWSOrganizations []entity.UserWorkspaceOrganization

	//
	userWSIds := h.GetParentIds(leaderIds)
	if len(userWSIds) == 0 {
		return nil
	}
	convertUserWSIds, err := utils.ConvertSliceStringToUUID(userWSIds)
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	userWS, err := h.postgresRepo.UserWorkspaceRepo.FindByFilter(ctx, &repository.FindUserWorkspaceByFilter{
		IDs: convertUserWSIds,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(userWS) != len(userWSIds) {
		return errors.New(errors.ErrCodeUserWorkspaceNotFound)
	}

	// Get user workspace organization data
	if tx != nil {
		userWSOrganizations, err = h.postgresRepo.UserWorkspaceOrganizationRepo.FindByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
			Filter: &repository.UserWorkspaceOrganizationFilter{
				UserWorkspaceIDs:      convertUserWSIds,
				IsRequireOrganization: true,
			},
			LockOption: clause.LockingStrengthShare,
			Tx:         tx,
		})
	} else {
		userWSOrganizations, err = h.postgresRepo.UserWorkspaceOrganizationRepo.FindByFilter(ctx, &repository.UserWorkspaceOrganizationFilter{
			UserWorkspaceIDs:      convertUserWSIds,
			IsRequireOrganization: true,
		})
	}
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(userWSOrganizations) != len(userWSIds) {
		return errors.New(errors.ErrCodeUserWorkspaceNotInOrganization)
	}

	/*
		Check format of parentOrganizationIds. Rule:
		- ParentOrganizationIds level must be in ASC order
	*/
	level := entity.ORGANiZATION_LEVEL_ROOT - 1
	for _, userWsId := range convertUserWSIds {
		for _, userWSOrganization := range userWSOrganizations {
			if userWsId != userWSOrganization.UserWorkspaceID {
				continue
			}
			if userWSOrganization.Organization.Level <= level {
				return errors.New(errors.ErrCodeInvalidParentOrganizationIds)
			}

			level = userWSOrganization.Organization.Level
		}
	}

	return nil
}
