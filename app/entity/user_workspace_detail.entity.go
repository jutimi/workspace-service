package entity

import (
	"time"
	"workspace-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserWorkspaceDetail struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name     string    `json:"name" gorm:"type:varchar(100);not null"`
	FullName string    `json:"full_name" gorm:"type:varchar(100);not null"`
	NameSlug string    `json:"name_slug" gorm:"type:varchar(100);not null"`
	BaseUserWorkspace

	// Relation
}

func NewUserWorkspaceDetail() *UserWorkspaceDetail {
	return &UserWorkspaceDetail{
		BaseUserWorkspace: BaseUserWorkspace{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

func (e *UserWorkspaceDetail) TableName() string {
	return "user_workspace_details"
}

func (e *UserWorkspaceDetail) BeforeSave(tx *gorm.DB) error {
	if e.FullName != "" {
		e.NameSlug = utils.ConvertToSnakeCase(e.FullName)
	}

	return nil
}
