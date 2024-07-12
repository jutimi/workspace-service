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
	WorkspaceID  *uuid.UUID
	IsActive     *bool
	Email        *string
	PhoneNumber  *string
	IDs          []uuid.UUID
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	Name         *string

	// Include option
	IsIncludeDetail bool
}
