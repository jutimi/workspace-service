package entity

import (
	"time"
	"workspace-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	ID                    uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name                  string     `json:"name" gorm:"type:varchar(100);not null"`
	NameSlug              string     `json:"name_slug" gorm:"type:varchar(100);not null"`
	Level                 int        `json:"level" gorm:"type:int;not null"`
	ParentOrganizationIDs *string    `json:"parent_organization_ids" gorm:"type:text"` // All parent organization ids
	ParentOrganizationID  *uuid.UUID `json:"parent_organization_id" gorm:"type:uuid"`  // Last parent organization id
	BaseWorkspace

	// Relation
	Workspace Workspace `gorm:"foreignKey:workspace_id;references:id"`
}

func NewOrganization() *Organization {
	return &Organization{
		BaseWorkspace: BaseWorkspace{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
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
