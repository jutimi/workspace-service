package model

import "github.com/google/uuid"

type CreateUserWorkspaceRequest struct {
	Detail CreateUserWorkspaceDetailRequest `json:"detail" validate:"required"`
}
type CreateUserWorkspaceResponse struct{}

type UpdateUserWorkspaceRequest struct {
	UserWorkspaceID uuid.UUID `json:"user_workspace_id" validate:"required"`
	CreateUserWorkspaceRequest
}
type UpdateUserWorkspaceResponse struct{}
