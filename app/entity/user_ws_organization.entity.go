package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ORGANiZATION_ROLE_LEADER     = "leader"
	ORGANiZATION_ROLE_SUB_LEADER = "sub_leader"
	ORGANiZATION_ROLE_MEMBER     = "member"
)

type UserWorkspaceOrganization struct {
	gorm.Model
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null"`
	Role           string    `json:"role" gorm:"type:varchar(20);not null"`
	LeaderIDs      *string   `json:"leader_ids" gorm:"type:text"`
	BaseUserWorkspace

	// Relation
	Organization Organization `gorm:"foreignKey:organization_id;references:id"`
}

func NewUserWorkspaceOrganization() *UserWorkspaceOrganization {
	return &UserWorkspaceOrganization{
		BaseUserWorkspace: BaseUserWorkspace{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

func (e *UserWorkspaceOrganization) TableName() string {
	return "user_workspace_organizations"
}
