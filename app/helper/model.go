package helper

import (
	"workspace-server/app/entity"
	"workspace-server/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateUserWsParams struct {
	Tx          *gorm.DB
	UserID      uuid.UUID
	WorkspaceID uuid.UUID
	PhoneNumber *string
	Email       *string
	Role        string
}

type CreateOrganizationParams struct {
	Tx                   *gorm.DB
	WorkspaceID          uuid.UUID
	ParentOrganizationId *uuid.UUID // Parent organization id (organization id)
	ParentLeaderId       *uuid.UUID // Manager of organization leader (user workspace id)
	Name                 string
	LeaderID             *uuid.UUID // Leader of organization (user workspace id)
	SubLeaders           []model.SubLeaderData
}

type CreateUserWorkspaceOrganizationParams struct {
	Tx               *gorm.DB
	Organization     *entity.Organization
	UserWorkspaceIds []uuid.UUID // Member ids (user workspace ids)
	LeaderId         *uuid.UUID  // Leader of members in current organization (user workspace id)
	LeaderIds        string      // List ids of mangers of leader (user workspace ids)
	Role             string
}

type UpdateOrganizationParams struct {
	Organization         *entity.Organization
	Tx                   *gorm.DB
	ParentOrganizationId *uuid.UUID // Parent organization id (organization id)
	ParentLeaderId       *uuid.UUID // Manager of organization leader (user workspace id)
	Name                 string
	LeaderID             *uuid.UUID // Leader of organization (user workspace id)
	SubLeaders           []model.SubLeaderData
}

type validateOrganizationParams struct {
	ParentOrganizationId *uuid.UUID
	ParentLeaderId       *uuid.UUID
}
type validateOrganizationResult struct {
	Organization              *entity.Organization
	UserWorkspaceOrganization *entity.UserWorkspaceOrganization
}
