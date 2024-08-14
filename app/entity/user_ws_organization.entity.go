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
	Id              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrganizationId  uuid.UUID `json:"organization_id" gorm:"type:uuid;not null"`
	Role            string    `json:"role" gorm:"type:varchar(20);not null"`
	LeaderIds       *string   `json:"leader_ids" gorm:"type:text"` // List ids of leader of user (user workspace id)
	UserWorkspaceId uuid.UUID `json:"user_workspace_id" gorm:"type:uuid;not null"`
	WorkspaceId     uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null"`
	CreatedAt       int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`

	// Relation
	Organization Organization `gorm:"foreignKey:organization_id;references:id"`
}

func NewUserWorkspaceOrganization() *UserWorkspaceOrganization {
	return &UserWorkspaceOrganization{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (e *UserWorkspaceOrganization) TableName() string {
	return "user_workspace_organizations"
}
