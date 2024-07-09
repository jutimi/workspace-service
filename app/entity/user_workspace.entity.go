package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ROLE_OWNER = "owner"
	ROLE_USER  = "user"
	ROLE_ADMIN = "admin"
)

type UserWorkspace struct {
	gorm.Model
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PhoneNumber *string   `json:"phone_number" gorm:"type:varchar(20);"`
	Email       *string   `json:"email" gorm:"type:varchar(100);"`
	Role        string    `json:"role" gorm:"type:enum('admin','user', 'owner');"`
	IsActive    bool      `json:"is_active" gorm:"default:true;type:bool;not null"`
	BaseWorkspaceEntity

	// Relation
	Workspace           Workspace           `gorm:"foreignKey:workspace_id;references:id"`
	UserWorkspaceDetail UserWorkspaceDetail `gorm:"foreignKey:user_workspace_id"`
}

func NewUserWorkspace() *UserWorkspace {
	return &UserWorkspace{
		BaseWorkspaceEntity: BaseWorkspaceEntity{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

func (e *UserWorkspace) TableName() string {
	return "user_workspaces"
}
