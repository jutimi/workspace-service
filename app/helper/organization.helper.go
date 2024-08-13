package helper

import (
	"context"
	"fmt"
	"strings"

	"workspace-server/app/entity"
	"workspace-server/app/model"
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

func (h *organizationHelper) CreateOrganization(
	ctx context.Context,
	data *CreateOrganizationParams,
) error {
	leaderIds := ""

	// Validate duplicate user workspace
	if err := h.validateDuplicateUserWorkspace(ctx, data.LeaderID, data.SubLeaders); err != nil {
		return err
	}

	// Create organization
	organization := entity.NewOrganization()
	organization.Name = data.Name
	organization.BaseWorkspace.WorkspaceID = data.WorkspaceID
	organization.Level = entity.ORGANiZATION_LEVEL_ROOT
	// Check parent organization
	parentOrganizationData, err := h.validateParentOrganization(ctx, &validateOrganizationParams{
		ParentOrganizationId: data.ParentOrganizationId,
		ParentLeaderId:       data.ParentLeaderId,
	})
	if err != nil {
		return err
	}
	if parentOrganizationData != nil {
		// Update organization base on parent organization data: Level, ParentOrganizationID, ParentOrganizationIDs, ManagerID
		organization.Level = parentOrganizationData.Organization.Level + 1
		organization.ParentOrganizationID = &parentOrganizationData.Organization.ID
		parentOrganizationIds := h.generateParentIds(
			ctx,
			*parentOrganizationData.Organization.ParentOrganizationIDs,
			parentOrganizationData.Organization.ID,
		)
		organization.ParentOrganizationIDs = &parentOrganizationIds
		organization.ManagerID = &parentOrganizationData.UserWorkspaceOrganization.ID

		// Update leader ids list
		leaderIds = h.generateParentIds(
			ctx,
			*parentOrganizationData.UserWorkspaceOrganization.LeaderIDs,
			parentOrganizationData.UserWorkspaceOrganization.UserWorkspaceID,
		)
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
			ctx,
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

func (h *organizationHelper) UpdateOrganization(
	ctx context.Context,
	data *UpdateOrganizationParams,
) error {
	leaderIds := ""

	if data.Organization == nil {
		return nil
	}

	// Validate duplicate user workspace
	if err := h.validateDuplicateUserWorkspace(ctx, data.LeaderID, data.SubLeaders); err != nil {
		return err
	}
	// Validate organization before update
	if err := h.validateUpdateOrganization(ctx, data); err != nil {
		return err
	}

	// Update organization
	// Check parent organization
	parentOrganizationData, err := h.validateParentOrganization(ctx, &validateOrganizationParams{
		ParentOrganizationId: data.ParentOrganizationId,
		ParentLeaderId:       data.ParentLeaderId,
	})
	if err != nil {
		return err
	}
	data.Organization.Name = data.Name
	if parentOrganizationData != nil {
		// Update organization base on parent organization data: Level, ParentOrganizationID, ParentOrganizationIDs, ManagerID
		data.Organization.Level = parentOrganizationData.Organization.Level + 1
		data.Organization.ParentOrganizationID = &parentOrganizationData.Organization.ID
		parentOrganizationIds := h.generateParentIds(
			ctx,
			*parentOrganizationData.Organization.ParentOrganizationIDs,
			parentOrganizationData.Organization.ID,
		)
		data.Organization.ParentOrganizationIDs = &parentOrganizationIds
		data.Organization.ManagerID = &parentOrganizationData.UserWorkspaceOrganization.ID

		// Update leader ids list
		leaderIds = h.generateParentIds(
			ctx,
			*parentOrganizationData.UserWorkspaceOrganization.LeaderIDs,
			parentOrganizationData.UserWorkspaceOrganization.UserWorkspaceID,
		)
	}
	if err := h.postgresRepo.OrganizationRepo.Update(ctx, data.Tx, data.Organization); err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	// Update user workspace organization
	// Remove all user workspace organizations of current organization
	if err := h.postgresRepo.UserWorkspaceOrganizationRepo.DeleteByFilter(ctx, data.Tx, &repository.FindUserWorkspaceOrganizationFilter{
		OrganizationID: &data.Organization.ID,
	}); err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	// Recreate user workspace organization
	if data.LeaderID != nil {
		if err := h.createUserWorkspaceOrganization(
			ctx,
			&CreateUserWorkspaceOrganizationParams{
				Organization: data.Organization,
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
			ctx,
			leaderIds,
			*data.LeaderID,
		)
	}
	if data.SubLeaders != nil {
		for _, subLeaderData := range data.SubLeaders {
			if err := h.createUserWorkspaceOrganization(ctx, &CreateUserWorkspaceOrganizationParams{
				Tx:               data.Tx,
				Organization:     data.Organization,
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

// ----------------------------------------------- Helper -----------------------------------------------

/*
Generate parent ids string (ex: 1/2/3/4)

- parentIds: previous parent ids (ex: 1/2/3)

- parentId: current parent id (ex: 4)
*/
func (h *organizationHelper) generateParentIds(
	ctx context.Context,
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
	limit := 1
	offset := 0

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
		Limit:            &limit,
		Offset:           &offset,
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
		leaderIds = h.generateParentIds(ctx, data.LeaderIds, *data.LeaderId)
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

/*
Function validate parent organization. Rule:

- Check exist both parent organization id and parent leader id

- Check exist parent organization

- Check exist user workspace

- Check user workspace is in parent organization
*/
func (h *organizationHelper) validateParentOrganization(
	ctx context.Context,
	data *validateOrganizationParams,
) (*validateOrganizationResult, error) {
	if (data.ParentOrganizationId != nil && data.ParentLeaderId == nil) ||
		(data.ParentOrganizationId == nil && data.ParentLeaderId != nil) {
		return nil, errors.New(errors.ErrCodeInvalidParentOrganizationData)
	}

	if data.ParentOrganizationId == nil && data.ParentLeaderId == nil {
		return nil, nil
	}

	// Get parent organization
	organization, err := h.postgresRepo.OrganizationRepo.FindOneByFilter(ctx, &repository.FindOrganizationByFilter{
		ID: data.ParentOrganizationId,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeOrganizationNotFound)
	}

	// Get leader data
	isActive := true
	leader, err := h.postgresRepo.UserWorkspaceRepo.FindOneByFilter(ctx, &repository.FindUserWorkspaceByFilter{
		ID:       data.ParentLeaderId,
		IsActive: &isActive,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserWorkspaceNotFound)
	}
	leaderOrganization, err := h.postgresRepo.UserWorkspaceOrganizationRepo.FindOneByFilter(ctx, &repository.FindUserWorkspaceOrganizationFilter{
		UserWorkspaceID: &leader.ID,
		OrganizationID:  &organization.ID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserWorkspaceNotInOrganization)
	}

	return &validateOrganizationResult{
		Organization:              organization,
		UserWorkspaceOrganization: leaderOrganization,
	}, nil
}

/*
Function validate organization before update. Rule:

- If organization has child, can only change organization name

- If organization doesn't have child, can change any data
*/
func (h *organizationHelper) validateUpdateOrganization(
	ctx context.Context,
	data *UpdateOrganizationParams,
) error {
	limit := 1
	offset := 0
	// Check existed child organization
	existedChildOrganization, err := h.postgresRepo.OrganizationRepo.FindByFilter(ctx, &repository.FindOrganizationByFilter{
		ParentOrganizationID: &data.Organization.ID,
		Limit:                &limit,
		Offset:               &offset,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedChildOrganization) == 0 {
		return nil
	}

	// Check update parent organization or manager
	if *data.ParentLeaderId != *data.Organization.ManagerID {
		return errors.New(errors.ErrCodeOrganizationHasChild)
	}
	if *data.ParentOrganizationId != *data.Organization.ParentOrganizationID {
		return errors.New(errors.ErrCodeOrganizationHasChild)
	}

	// Check update user workspace organization
	userWorkspaceIds := make([]uuid.UUID, 0)
	if data.LeaderID != nil {
		userWorkspaceIds = append(userWorkspaceIds, *data.LeaderID)
	}
	for _, subLeader := range data.SubLeaders {
		userWorkspaceIds = append(userWorkspaceIds, subLeader.SubLeaderId)
		userWorkspaceIds = append(userWorkspaceIds, subLeader.MemberIds...)
	}

	userWSOrganizations, err := h.postgresRepo.UserWorkspaceOrganizationRepo.FindByFilter(ctx, &repository.FindUserWorkspaceOrganizationFilter{
		OrganizationID: &data.Organization.ID,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(userWSOrganizations) != len(userWorkspaceIds) {
		return errors.New(errors.ErrCodeOrganizationHasChild)
	}

	for _, userWSOrganization := range userWSOrganizations {
		if !utils.IsSliceContain(userWSOrganization.UserWorkspaceID, userWorkspaceIds) {
			return errors.New(errors.ErrCodeOrganizationHasChild)
		}
	}

	return nil
}

/*
Function validate duplicate user workspace. Rule:

- Check user input have duplicate user workspace in 1 organization
*/
func (h *organizationHelper) validateDuplicateUserWorkspace(
	ctx context.Context,
	leaderId *uuid.UUID,
	subLeaders []model.SubLeaderData,
) error {
	userWorkspaceIds := make([]uuid.UUID, 0)
	if leaderId != nil {
		userWorkspaceIds = append(userWorkspaceIds, *leaderId)
	}
	for _, subLeader := range subLeaders {
		userWorkspaceIds = append(userWorkspaceIds, subLeader.SubLeaderId)
		userWorkspaceIds = append(userWorkspaceIds, subLeader.MemberIds...)
	}

	userWorkspaceMemo := make(map[uuid.UUID]uuid.UUID, 0)
	for _, userWorkspaceId := range userWorkspaceIds {
		if _, ok := userWorkspaceMemo[userWorkspaceId]; ok {
			return errors.New(errors.ErrCodeDuplicateUserWorkspaceInOrganization)
		}

		userWorkspaceMemo[userWorkspaceId] = userWorkspaceId
	}

	return nil
}
