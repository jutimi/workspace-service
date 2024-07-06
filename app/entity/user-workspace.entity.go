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
	WorkspaceID uuid.UUID `json:"workspace_id" gorm:"type:uuid;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func NewUserWorkspace() *UserWorkspace {
	return &UserWorkspace{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (e *UserWorkspace) TableName() string {
	return "user_workspaces"
}
