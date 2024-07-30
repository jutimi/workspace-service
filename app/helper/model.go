package helper

import (
	"workspace-server/app/entity"

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
	ParentOrganizationId *uuid.UUID // Parent organization id
	ParentLeaderId       *uuid.UUID // Leader of organization leader
	Name                 string
	Leader               *uuid.UUID
	SubLeaders           []SubLeaderData
}
type SubLeaderData struct {
	SubLeaderId uuid.UUID
	MemberIds   []uuid.UUID
}

type CreateUserWorkspaceOrganizationParams struct {
	Tx               *gorm.DB
	Organization     *entity.Organization
	UserWorkspaceIds []uuid.UUID
	LeaderIds        []uuid.UUID
	Role             string
}
