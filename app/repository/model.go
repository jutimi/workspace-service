package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FindByFilterForUpdateParams struct {
	Filter     interface{}
	LockOption string
	Tx         *gorm.DB
}

type FindWorkspaceByFilter struct {
	ID           *uuid.UUID
	IsActive     *bool
	Email        *string
	PhoneNumber  *string
	IDs          []uuid.UUID
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	Name         *string
}

type FindUserWorkspaceByFilter struct {
	ID           *uuid.UUID
	IDs          []uuid.UUID
	WorkspaceID  *uuid.UUID
	WorkspaceIDs []uuid.UUID
	UserID       *uuid.UUID
	UserIDs      []uuid.UUID
	IsActive     *bool
	Email        *string
	PhoneNumber  *string
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	Name         *string
	Role         *string

	// Include option
	IsIncludeDetail bool // Left join user workspace detail

	// Require option
	IsRequireDetail bool // Inner join user workspace detail
}

type FindOrganizationByFilter struct {
	ID                   *uuid.UUID
	IDs                  []uuid.UUID
	WorkspaceID          *uuid.UUID
	ParentOrganizationID *uuid.UUID
	Name                 *string
	Limit                *int
	Offset               *int
	Level                *int
}

type FindUserWorkspaceOrganizationFilter struct {
	ID               *uuid.UUID
	IDs              []uuid.UUID
	UserWorkspaceID  *uuid.UUID
	UserWorkspaceIDs []uuid.UUID
	WorkspaceID      *uuid.UUID
	WorkspaceIDs     []uuid.UUID
	OrganizationID   *uuid.UUID
	OrganizationIDs  []uuid.UUID
	Role             *string
	Limit            *int
	Offset           *int

	// Include option
	IsIncludeOrganization bool // Left join organization

	// Require option
	IsRequireOrganization bool // Inner join organization
}
