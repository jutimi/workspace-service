package entity

import (
	"time"

	"workspace-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserWorkspaceDetail struct {
	gorm.Model
	Id              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name            string    `json:"name" gorm:"type:varchar(100);not null"`
	FullName        string    `json:"full_name" gorm:"type:varchar(100);not null"`
	NameSlug        string    `json:"full_name_slug" gorm:"type:varchar(100);not null"`
	UserWorkspaceId uuid.UUID `json:"user_workspace_id" gorm:"type:uuid;not null"`
	WorkspaceId     uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null"`
	CreatedAt       int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`

	// Relation
}

func NewUserWorkspaceDetail() *UserWorkspaceDetail {
	return &UserWorkspaceDetail{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (e *UserWorkspaceDetail) TableName() string {
	return "user_workspace_details"
}

func (e *UserWorkspaceDetail) BeforeSave(tx *gorm.DB) error {
	if e.FullName != "" {
		e.NameSlug = utils.Slugify(e.FullName)
	}

	return nil
}
