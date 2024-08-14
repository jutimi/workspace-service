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
	Id           *uuid.UUID
	IsActive     *bool
	Email        *string
	PhoneNumber  *string
	Ids          []uuid.UUID
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	Name         *string
}

type FindUserWorkspaceByFilter struct {
	Id           *uuid.UUID
	Ids          []uuid.UUID
	WorkspaceId  *uuid.UUID
	WorkspaceIds []uuid.UUID
	UserId       *uuid.UUID
	UserIds      []uuid.UUID
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
	Id                       *uuid.UUID
	Ids                      []uuid.UUID
	WorkspaceId              *uuid.UUID
	ParentOrganizationId     *uuid.UUID
	Name                     *string
	Limit                    *int
	Offset                   *int
	Level                    *int
	ParentOrganizationIdsStr *string
}

type FindUserWorkspaceOrganizationFilter struct {
	Id               *uuid.UUID
	Ids              []uuid.UUID
	UserWorkspaceId  *uuid.UUID
	UserWorkspaceIds []uuid.UUID
	WorkspaceId      *uuid.UUID
	WorkspaceIds     []uuid.UUID
	OrganizationId   *uuid.UUID
	OrganizationIds  []uuid.UUID
	Role             *string
	Limit            *int
	Offset           *int
	LeaderIds        *string

	// Include option
	IsIncludeOrganization bool // Left join organization

	// Require option
	IsRequireOrganization bool // Inner join organization
}
