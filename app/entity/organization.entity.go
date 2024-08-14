package entity

import (
	"time"

	"workspace-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ORGANiZATION_LEVEL_ROOT = 0
)

type Organization struct {
	gorm.Model
	ID                    uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name                  string     `json:"name" gorm:"type:varchar(100);not null"`
	NameSlug              string     `json:"name_slug" gorm:"type:varchar(100);not null"`
	Level                 int        `json:"level" gorm:"type:int;not null"`
	ParentOrganizationIDs *string    `json:"parent_organization_ids" gorm:"type:text"` // List ids of parent organization
	ParentOrganizationID  *uuid.UUID `json:"parent_organization_id" gorm:"type:uuid"`  // Current parent organization
	ManagerID             *uuid.UUID `json:"manager_id" gorm:"type:uuid"`              // Manager of leader of organization (user workspace id)
	WorkspaceID           uuid.UUID  `json:"workspace_id" gorm:"type:uuid;not null"`
	CreatedAt             int64      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             int64      `json:"updated_at" gorm:"autoUpdateTime:milli"`

	// Relation
	Workspace Workspace `gorm:"foreignKey:workspace_id;references:id"`
}

func NewOrganization() *Organization {
	return &Organization{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (e *Organization) TableName() string {
	return "organizations"
}

func (e *Organization) BeforeSave(tx *gorm.DB) error {
	if e.Name != "" {
		e.NameSlug = utils.Slugify(e.Name)
	}

	return nil
}
