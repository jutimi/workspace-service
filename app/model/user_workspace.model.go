package model

import "github.com/google/uuid"

type CreateUserWorkspaceRequest struct {
	Name        string                           `json:"name" validate:"required"`
	Email       *string                          `json:"email"`
	PhoneNumber *string                          `json:"phone_number"`
	Detail      CreateUserWorkspaceDetailRequest `json:"detail" validate:"required"`
}
type CreateUserWorkspaceResponse struct{}

type UpdateUserWorkspaceRequest struct {
	UserWorkspaceID uuid.UUID `json:"user_workspace_id" validate:"required"`
	CreateUserWorkspaceRequest
}
type UpdateUserWorkspaceResponse struct{}

type InactiveUserWorkspaceRequest struct{}
type InactiveUserWorkspaceResponse struct{}

type RemoveUserWorkspaceRequest struct{}
type RemoveUserWorkspaceResponse struct{}
