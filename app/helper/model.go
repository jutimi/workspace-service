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
	Tx                 *gorm.DB
	WorkspaceID        uuid.UUID
	ParentOrganization *entity.Organization
	Name               string
	Leader             *MemberInfo
	SubLeaders         []MemberInfo
	Members            []MemberInfo
}
type MemberInfo struct {
	LeaderIds            *string // List of leader ids of user
	entity.UserWorkspace         // User workspace info of user
}

type CreateUserWorkspaceOrganizationParams struct {
	Tx           *gorm.DB
	Organization *entity.Organization
	Data         []MemberInfo
	Role         string
}
