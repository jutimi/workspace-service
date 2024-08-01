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
	leaderIds := ""
	isActive := true

	// Create organization
	organization := entity.NewOrganization()
	organization.Name = data.Name
	organization.BaseWorkspace.WorkspaceID = data.WorkspaceID
	organization.Level = entity.ORGANiZATION_LEVEL_ROOT
	// Check parent organization
	if data.ParentOrganizationId != nil && data.ParentLeaderId != nil {
		// Get parent organization
		parentOrganization, err := h.postgresRepo.OrganizationRepo.FindOneByFilter(ctx, &repository.FindOrganizationByFilter{
			ID: data.ParentOrganizationId,
		})
		if err != nil {
			return errors.New(errors.ErrCodeOrganizationNotFound)
		}

		// Update organization base on parent organization data: Level, ParentOrganizationID, ParentOrganizationIDs
		organization.Level = parentOrganization.Level + 1
		organization.ParentOrganizationID = &parentOrganization.ID
		parentOrganizationIds := h.generateParentIds(
			*parentOrganization.ParentOrganizationIDs,
			parentOrganization.ID,
		)
		organization.ParentOrganizationIDs = &parentOrganizationIds

		// Get leader data
		leader, err := h.postgresRepo.UserWorkspaceRepo.FindOneByFilter(ctx, &repository.FindUserWorkspaceByFilter{
			ID:       data.ParentLeaderId,
			IsActive: &isActive,
		})
		if err != nil {
			return errors.New(errors.ErrCodeUserWorkspaceNotFound)
		}
		leaderOrganization, err := h.postgresRepo.UserWorkspaceOrganizationRepo.FindOneByFilter(ctx, &repository.FindUserWorkspaceOrganizationFilter{
			UserWorkspaceID: &leader.ID,
			OrganizationID:  &parentOrganization.ID,
		})
		if err != nil {
			return errors.New(errors.ErrCodeUserWorkspaceNotInOrganization)
		}

		// Update organization base on parent organization data: ParentOrganizationLeaderID
		organization.ManagerID = &leader.ID

		// Update leader ids list
		leaderIds = h.generateParentIds(
			*leaderOrganization.LeaderIDs,
			leader.ID,
		)
	} else if (data.ParentOrganizationId != nil && data.ParentLeaderId == nil) ||
		(data.ParentOrganizationId == nil && data.ParentLeaderId != nil) {
		return errors.New(errors.ErrCodeInvalidParentOrganizationData)
	}
	if err := h.postgresRepo.OrganizationRepo.Create(ctx, data.Tx, organization); err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	// Create user workspace organization
	if data.LeaderID != nil {
		if err := h.createUserWorkspaceOrganization(
			ctx,
			&CreateUserWorkspaceOrganizationParams{
				Organization: organization,
				Role:         entity.ORGANiZATION_ROLE_LEADER,
				UserWorkspaceIds: []uuid.UUID{
					*data.LeaderID,
				},
				LeaderId:  nil,
				LeaderIds: leaderIds,
				Tx:        data.Tx,
			},
		); err != nil {
			return errors.New(errors.ErrCodeInternalServerError)
		}

		// Update leader ids list
		leaderIds = h.generateParentIds(
			leaderIds,
			*data.LeaderID,
		)
	}
	if data.SubLeaders != nil {
		for _, subLeaderData := range data.SubLeaders {
			if err := h.createUserWorkspaceOrganization(ctx, &CreateUserWorkspaceOrganizationParams{
				Tx:               data.Tx,
				Organization:     organization,
				UserWorkspaceIds: subLeaderData.MemberIds,
				LeaderId:         &subLeaderData.SubLeaderId,
				LeaderIds:        leaderIds,
				Role:             entity.ORGANiZATION_ROLE_SUB_LEADER,
			}); err != nil {
				return errors.New(errors.ErrCodeInternalServerError)
			}
		}
	}

	return nil
}

func (h *organizationHelper) GetParentIds(parentIdStr string) []uuid.UUID {
	parentIds := strings.Split(parentIdStr, "/")
	result, _ := utils.ConvertSliceStringToUUID(parentIds)
	return result
}

// ----------------------------------------------- Helper -----------------------------------------------

/*
Generate parent ids string (ex: 1/2/3/4)

- parentIds: previous parent ids (ex: 1/2/3)

- parentId: current parent id (ex: 4)
*/
func (h *organizationHelper) generateParentIds(
	parentIds string,
	parentId uuid.UUID,
) string {
	return fmt.Sprintf("%s/%s", parentIds, parentId.String())
}

func (h *organizationHelper) createUserWorkspaceOrganization(
	ctx context.Context,
	data *CreateUserWorkspaceOrganizationParams,
) error {
	isActive := true

	if data.Organization == nil {
		return nil
	}

	// Check user workspace data
	userWS, err := h.postgresRepo.UserWorkspaceRepo.FindByFilter(ctx, &repository.FindUserWorkspaceByFilter{
		IDs:      data.UserWorkspaceIds,
		IsActive: &isActive,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(userWS) != len(data.UserWorkspaceIds) {
		return errors.New(errors.ErrCodeUserWorkspaceNotFound)
	}

	// Check user workspace in parent organization
	parentOrganizationIds := h.GetParentIds(*data.Organization.ParentOrganizationIDs)
	existedUserWSOrganization, err := h.postgresRepo.UserWorkspaceOrganizationRepo.FindByFilter(ctx, &repository.FindUserWorkspaceOrganizationFilter{
		UserWorkspaceIDs: data.UserWorkspaceIds,
		OrganizationIDs:  parentOrganizationIds,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUserWSOrganization) > 0 {
		return errors.New(errors.ErrCodeUserWorkspaceAlreadyInOrganization)
	}

	// Create user workspace organization
	leaderIds := data.LeaderIds
	if data.LeaderId != nil {
		leaderIds = h.generateParentIds(data.LeaderIds, *data.LeaderId)
	}
	userWSOrganizations := make([]entity.UserWorkspaceOrganization, 0)
	for _, userWSId := range data.UserWorkspaceIds {
		userWorkspaceOrganization := entity.NewUserWorkspaceOrganization()
		userWorkspaceOrganization.OrganizationID = data.Organization.ID
		userWorkspaceOrganization.BaseUserWorkspace.WorkspaceID = data.Organization.WorkspaceID
		userWorkspaceOrganization.BaseUserWorkspace.UserWorkspaceID = userWSId
		userWorkspaceOrganization.Role = data.Role
		userWorkspaceOrganization.LeaderIDs = &leaderIds
		userWSOrganizations = append(userWSOrganizations, *userWorkspaceOrganization)
	}
	if err := h.postgresRepo.UserWorkspaceOrganizationRepo.BulkCreate(ctx, data.Tx, userWSOrganizations); err != nil {
		return err
	}

	return nil
}

func (h *organizationHelper) validateParentOrganizationIds(
	ctx context.Context,
	parentOrganizationIds string,
) error {
	var organizations []entity.Organization

	//
	organizationIds := h.GetParentIds(parentOrganizationIds)
	if len(organizationIds) == 0 {
		return nil
	}

	// Get organizations
	organizations, err := h.postgresRepo.OrganizationRepo.FindByFilter(ctx, &repository.FindOrganizationByFilter{
		IDs: organizationIds,
	})
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
	for _, organizationId := range organizationIds {
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
