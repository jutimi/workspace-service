package helper

import (
	"context"
	"fmt"
	"strings"
	"workspace-server/app/entity"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/package/errors"

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
	// Create organization
	organization := entity.NewOrganization()
	organization.Name = data.Name
	organization.BaseWorkspace.WorkspaceID = data.WorkspaceID
	organization.Level = 0
	if data.ParentOrganization != nil {
		organization.Level = data.ParentOrganization.Level + 1
		organization.ParentOrganizationID = &data.ParentOrganization.ID
		parentOrganizationIds := h.generateParentOrganizationIds(*data.ParentOrganization.ParentOrganizationIDs, data.ParentOrganization.ID)
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

func (h *organizationHelper) GetParentOrganizationIds(organizationIDs string) []string {
	return strings.Split(organizationIDs, "/")
}

// ----------------------------------------------- Helper -----------------------------------------------
func (h *organizationHelper) generateParentOrganizationIds(parentOrganizationIds string, parentOrganizationId uuid.UUID) string {
	return fmt.Sprintf("%s/%s", parentOrganizationIds, parentOrganizationId.String())
}

func (h *organizationHelper) createUserWorkspaceOrganization(ctx context.Context, data *CreateUserWorkspaceOrganizationParams) error {
	if data.Organization == nil {
		return nil
	}

	users := make([]entity.UserWorkspaceOrganization, 0)
	for _, member := range data.Data {
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

func (h *organizationHelper) validateParentOrganizationIds(parentOrganizationIds string) error {
	return nil
}

func (h *organizationHelper) validateLeaderIds(leaderIds string) error {
	return nil
}
