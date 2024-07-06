package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PhoneNumber *string   `json:"phone_number" gorm:"type:varchar(20);"`
	Address     *string   `json:"address" gorm:"type:text;"`
	Email       *string   `json:"email" gorm:"type:varchar(100);"`
	IsActive    bool      `json:"is_active" gorm:"default:true;type:bool;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func NewWorkspace() *Workspace {
	return &Workspace{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (e *Workspace) TableName() string {
	return "workspaces"
}
