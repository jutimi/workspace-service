package entity

import "github.com/google/uuid"

type BaseWorkspace struct {
	WorkspaceID uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type BaseUserWorkspace struct {
	UserWorkspaceID uuid.UUID `json:"user_workspace_id" gorm:"type:uuid;not null"`
	WorkspaceID     uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null"`
	CreatedAt       int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
