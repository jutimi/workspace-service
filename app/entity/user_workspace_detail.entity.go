package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserWorkspaceDetail struct {
	gorm.Model
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BaseUserWorkspaceEntity

	// Relation
}

func NewUserWorkspaceDetail() *UserWorkspaceDetail {
	return &UserWorkspaceDetail{
		BaseUserWorkspaceEntity: BaseUserWorkspaceEntity{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

func (e *UserWorkspaceDetail) TableName() string {
	return "user_workspace_details"
}
