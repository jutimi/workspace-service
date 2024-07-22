package model

import "github.com/google/uuid"

type CreateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address"`
	Email       string `json:"email" validate:"required"`
}
type CreateWorkspaceResponse struct{}

type UpdateWorkspaceRequest struct {
	ID uuid.UUID `query:"id" validate:"required"`
	CreateWorkspaceRequest
}
type UpdateWorkspaceResponse struct{}
