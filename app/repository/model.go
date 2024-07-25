package repository

import "github.com/google/uuid"

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

	// Include option
	IsIncludeDetail bool // Left join user workspace detail

	// Require option
	IsRequireDetail bool // Inner join user workspace detail
}
